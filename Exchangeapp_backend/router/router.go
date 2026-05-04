package router

import (
	"exchangeapp/internal/handler"
	"exchangeapp/internal/middleware"
	"exchangeapp/pkg/auth"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
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

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Health check
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

	// Auth routes (public)
	auth := r.Group("/api/auth")
	{
		auth.POST("/login", h.Auth.Login)
		auth.POST("/register", h.Auth.Register)
	}

	// API routes
	api := r.Group("/api")
	api.GET("/exchangeRates", h.Exchange.GetAll)
	api.Use(middleware.AuthMiddleware(jwt))
	{
		api.POST("/exchangeRates", h.Exchange.Create)
		api.POST("/articles", h.Article.Create)
		api.GET("/articles", h.Article.GetAll)
		api.GET("/articles/:id", h.Article.GetByID)
		api.POST("/articles/:id/like", h.Like.Like)
		api.GET("/articles/:id/like", h.Like.GetLikes)
	}

	return r
}
