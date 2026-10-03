package router

import (
	"currencyflow/internal/handler"
	"currencyflow/internal/middleware"
	"currencyflow/internal/repository"
	"currencyflow/pkg/auth"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"gorm.io/gorm"
)

type Handlers struct {
	Auth         *handler.AuthHandler
	Article      *handler.ArticleHandler
	Exchange     *handler.ExchangeHandler
	Like         *handler.LikeHandler
	RateHistory  *handler.RateHistoryHandler
	WS           *handler.WSHandler
	Favorite     *handler.FavoriteHandler
	Alert        *handler.AlertHandler
	Notification *handler.NotificationHandler
	Post         *handler.PostHandler
	Follow       *handler.FollowHandler
	UserProfile  *handler.UserProfileHandler
	AIAnalyst    *handler.AIAnalystHandler
}

// SetupRouter builds the gin.Engine and returns it along with the rate limiter
// (so the caller can call limiter.Stop() during graceful shutdown).
func SetupRouter(h Handlers, jwt *auth.JWTManager, db *gorm.DB, userRepo repository.UserRepository) (*gin.Engine, *middleware.IPRateLimiter, *middleware.IPRateLimiter) {
	r := gin.Default()
	if err := r.SetTrustedProxies(trustedProxiesFromEnv()); err != nil {
		log.Fatalf("invalid TRUSTED_PROXIES: %v", err)
	}

	// CORS origins from env (comma-separated), fallback to localhost
	corsOrigins := os.Getenv("CORS_ALLOWED_ORIGINS")
	if corsOrigins == "" {
		corsOrigins = "http://localhost:5173,http://localhost,http://localhost:80"
	}
	origins := strings.Split(corsOrigins, ",")
	for i := range origins {
		origins[i] = strings.TrimSpace(origins[i])
	}

	// Global middleware
	r.Use(middleware.Tracing())
	r.Use(middleware.MetricsAndLogging())
	r.Use(middleware.SecurityHeaders())
	rateLimiterFn, rateLimiter := middleware.RateLimit(100) // 100 req/s per IP
	r.Use(rateLimiterFn)
	r.Use(cors.New(cors.Config{
		AllowOrigins:     origins,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
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

	// Public auth routes (stricter rate limit: 10 req/s per IP)
	authRateLimitFn, authRateLimiter := middleware.RateLimit(10)
	authGroup := v1.Group("/auth")
	authGroup.Use(authRateLimitFn)
	{
		authGroup.POST("/login", h.Auth.Login)
		authGroup.POST("/register", h.Auth.Register)
	}

	// Public read routes
	v1.GET("/exchangeRates", h.Exchange.GetAll)
	v1.GET("/rates/history", h.RateHistory.GetHistory)
	v1.GET("/rates/latest", h.RateHistory.GetLatest)

	// AI Analyst (requires auth to prevent billing abuse)
	v1.POST("/ai/analyze", middleware.AuthMiddleware(jwt, userRepo), h.AIAnalyst.Analyze)

	// Public social routes
	v1.GET("/posts", h.Post.GetAll)
	v1.GET("/users/:id", h.UserProfile.GetProfile)

	// WebSocket
	v1.GET("/ws", h.WS.HandleWebSocket)
	v1.GET("/ws/clients", h.WS.GetClientCount)

	// Protected routes
	protected := v1.Group("")
	protected.Use(middleware.AuthMiddleware(jwt, userRepo))
	{
		protected.POST("/exchangeRates", h.Exchange.Create)
		protected.POST("/articles", h.Article.Create)
		protected.GET("/articles", h.Article.GetAll)
		protected.GET("/articles/:id", h.Article.GetByID)
		protected.POST("/articles/:id/like", h.Like.Like)
		protected.GET("/articles/:id/like", h.Like.GetLikes)

		// Favorites
		protected.POST("/favorites", h.Favorite.Add)
		protected.GET("/favorites", h.Favorite.GetByUser)
		protected.DELETE("/favorites", h.Favorite.Remove)
		protected.GET("/favorites/check", h.Favorite.Check)

		// Alerts
		protected.POST("/alerts", h.Alert.Create)
		protected.GET("/alerts", h.Alert.GetByUser)
		protected.DELETE("/alerts/:id", h.Alert.Delete)

		// Notifications
		protected.GET("/notifications", h.Notification.GetAll)
		protected.PUT("/notifications/:id/read", h.Notification.MarkRead)
		protected.PUT("/notifications/read-all", h.Notification.MarkAllRead)
		protected.GET("/notifications/unread-count", h.Notification.CountUnread)

		// Social
		protected.POST("/posts", h.Post.Create)
		protected.POST("/posts/:id/like", h.Post.Like)
		protected.POST("/users/:id/follow", h.Follow.Follow)
		protected.DELETE("/users/:id/follow", h.Follow.Unfollow)
		protected.GET("/users/:id/following", h.Follow.IsFollowing)
		protected.PUT("/users/profile", h.UserProfile.UpdateProfile)
	}

	// Deprecated: backward-compatible routes under /api/* — use /api/v1/* instead.
	// These will be removed in a future version.
	r.POST("/api/auth/login", h.Auth.Login)
	r.POST("/api/auth/register", h.Auth.Register)
	r.GET("/api/exchangeRates", h.Exchange.GetAll)

	return r, rateLimiter, authRateLimiter
}

// trustedProxiesFromEnv returns the CIDRs/addresses whose X-Forwarded-* headers
// may be trusted. Defaulting to loopback plus private container ranges keeps
// ClientIP() useful behind Docker/Nginx without trusting arbitrary clients.
func trustedProxiesFromEnv() []string {
	raw := strings.TrimSpace(os.Getenv("TRUSTED_PROXIES"))
	if raw == "none" {
		return nil
	}
	if raw == "" {
		raw = "127.0.0.1/32,::1/128,10.0.0.0/8,172.16.0.0/12,192.168.0.0/16"
	}
	parts := strings.Split(raw, ",")
	out := make([]string, 0, len(parts))
	for _, part := range parts {
		if item := strings.TrimSpace(part); item != "" {
			out = append(out, item)
		}
	}
	return out
}
