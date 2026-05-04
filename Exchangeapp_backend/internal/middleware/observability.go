package middleware

import (
	"exchangeapp/pkg/logger"
	"exchangeapp/pkg/metrics"
	"fmt"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

// Tracing 返回 OpenTelemetry 链路追踪中间件
func Tracing() gin.HandlerFunc {
	return otelgin.Middleware("exchangeapp")
}

// MetricsAndLogging 返回 Prometheus 指标 + 结构化日志中间件
func MetricsAndLogging() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.FullPath() // 使用路由模板而非实际路径，避免基数爆炸

		// 并发数 +1
		metrics.InflightRequests.Inc()
		defer metrics.InflightRequests.Dec()

		// 将 TraceID 注入到 gin.Context 供下游使用
		spanCtx := trace.SpanContextFromContext(c.Request.Context())
		if spanCtx.IsValid() {
			c.Set("trace_id", spanCtx.TraceID().String())
		}

		// 记录请求开始（结构化日志）
		reqLogger := logger.WithRequest(c.Request.Context(), c.Request.Method, path)
		reqLogger.Info("request started",
			zap.String("client_ip", c.ClientIP()),
			zap.String("user_agent", c.Request.UserAgent()),
		)

		// 执行后续 handler
		c.Next()

		// 计算延迟
		duration := time.Since(start).Seconds()
		status := strconv.Itoa(c.Writer.Status())

		// Prometheus 指标
		metrics.RequestsTotal.WithLabelValues(c.Request.Method, path, status).Inc()
		metrics.RequestDuration.WithLabelValues(c.Request.Method, path, status).Observe(duration)

		// 在当前 span 上附加业务属性
		span := trace.SpanFromContext(c.Request.Context())
		if span.IsRecording() {
			span.SetAttributes(
				attribute.Int("http.status_code", c.Writer.Status()),
				attribute.Float64("http.duration_ms", duration*1000),
			)
			if c.Writer.Status() >= 500 {
				span.SetStatus(codes.Error, fmt.Sprintf("HTTP %d", c.Writer.Status()))
			}
		}

		// 请求结束日志
		reqLogger.Info("request completed",
			zap.Int("status", c.Writer.Status()),
			zap.Float64("duration_ms", duration*1000),
			zap.Int("body_size", c.Writer.Size()),
		)
	}
}
