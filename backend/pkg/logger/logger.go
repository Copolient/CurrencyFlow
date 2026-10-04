package logger

import (
	"context"
	"os"

	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var log *zap.Logger

func Init() {
	cfg := zap.NewProductionConfig()
	cfg.OutputPaths = []string{"stdout"}
	cfg.EncoderConfig.TimeKey = "timestamp"
	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	cfg.EncoderConfig.StacktraceKey = "stacktrace"

	var err error
	log, err = cfg.Build()
	if err != nil {
		panic(err)
	}
}

func L() *zap.Logger {
	return log
}

// WithTraceID 从 context 中提取 TraceID 注入到 logger
func WithTraceID(ctx context.Context) *zap.Logger {
	spanCtx := trace.SpanContextFromContext(ctx)
	if spanCtx.IsValid() {
		return log.With(
			zap.String("trace_id", spanCtx.TraceID().String()),
			zap.String("span_id", spanCtx.SpanID().String()),
		)
	}
	return log
}

// WithRequest 附加请求级字段
func WithRequest(ctx context.Context, method, path string) *zap.Logger {
	return WithTraceID(ctx).With(
		zap.String("method", method),
		zap.String("path", path),
	)
}

func Sync() {
	if log != nil {
		_ = log.Sync()
	}
}

func init() {
	// fallback: 如果 Init() 未调用，提供一个基本 logger
	if log == nil {
		log = zap.New(zapcore.NewCore(
			zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
			zapcore.Lock(os.Stdout),
			zap.InfoLevel,
		))
	}
}
