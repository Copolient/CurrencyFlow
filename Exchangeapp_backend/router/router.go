package router

import (
	"exchangeapp/internal/handler"
	"exchangeapp/internal/middleware"
	"exchangeapp/pkg/auth"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"gorm.io/gorm"
)

type Handlers struct {
	Auth     *handler.AuthHandler
	Article  *handler.ArticleHandler
	Exchange *handler.ExchangeHandler
	Like     *handler.LikeHandler
}

func SetupRouter(h Handlers, jwt *auth.JWTManager, db *gorm.DB) *gin.Engine {
	r := gin.Default()

	// Global middleware
	r.Use(middleware.Tracing())
	r.Use(middleware.MetricsAndLogging())
	r.Use(middleware.SecurityHeaders())
	r.Use(middleware.RateLimit(100)) // 100 req/s per IP
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Health & Metrics (no versioning)
	r.GET("/healthz", func(c *gin.Context) {
		status := gin.H{"status": "ok"}
		sqlDB, err := db.DB()
		if err != nil || sqlDB.Ping() != nil {
			status["database"] = "unhealthy"
			c.JSON(http.StatusServiceUnavailable, status)
			return
		}
		status["database"] = "ok"
		c.JSON(http.StatusOK, status)
	})
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// API v1
	v1 := r.Group("/api/v1")

	// Public auth routes
	authGroup := v1.Group("/auth")
	{
		authGroup.POST("/login", h.Auth.Login)
		authGroup.POST("/register", h.Auth.Register)
	}

	// Public read routes
	v1.GET("/exchangeRates", h.Exchange.GetAll)

	// Protected routes
	protected := v1.Group("")
	protected.Use(middleware.AuthMiddleware(jwt))
	{
		protected.POST("/exchangeRates", h.Exchange.Create)
		protected.POST("/articles", h.Article.Create)
		protected.GET("/articles", h.Article.GetAll)
		protected.GET("/articles/:id", h.Article.GetByID)
		protected.POST("/articles/:id/like", h.Like.Like)
		protected.GET("/articles/:id/like", h.Like.GetLikes)
	}

	// Backward-compatible redirects from /api/* to /api/v1/*
	r.POST("/api/auth/login", h.Auth.Login)
	r.POST("/api/auth/register", h.Auth.Register)
	r.GET("/api/exchangeRates", h.Exchange.GetAll)

	return r
}
