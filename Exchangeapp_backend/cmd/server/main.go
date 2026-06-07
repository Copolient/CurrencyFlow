package main

import (
	"context"
	"exchangeapp/internal/handler"
	"exchangeapp/internal/repository"
	"exchangeapp/internal/scheduler"
	"exchangeapp/internal/service"
	ws "exchangeapp/internal/websocket"
	"exchangeapp/migrations"
	"exchangeapp/pkg/auth"
	"exchangeapp/pkg/cache"
	"exchangeapp/pkg/config"
	"exchangeapp/pkg/database"
	"exchangeapp/pkg/logger"
	"exchangeapp/pkg/tracing"
	"exchangeapp/router"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	migrateFlag := flag.Bool("migrate", false, "Run database migrations and exit")
	flag.Parse()

	// Load config
	cfg := config.Load()

	// Initialize structured logging
	logger.Init()
	defer logger.Sync()

	// Initialize OpenTelemetry tracing
	otelShutdown, err := tracing.Init("exchangeapp", os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT"))
	if err != nil {
		log.Printf("Warning: failed to init tracing: %v", err)
	} else {
		defer func() {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			_ = otelShutdown(ctx)
		}()
	}

	// Initialize infrastructure
	db := database.NewMySQL(cfg.Database)
	redisCache := cache.NewRedisCache(cfg.Redis)
	jwt := auth.NewJWTManager(cfg.JWT.Secret)

	// Run migrations if requested
	if *migrateFlag {
		migrations.Run(db)
		return
	}

	// Initialize repositories
	userRepo := repository.NewUserRepository(db)
	articleRepo := repository.NewArticleRepository(db)
	exchangeRepo := repository.NewExchangeRateRepository(db)
	rateHistoryRepo := repository.NewRateHistoryRepository(db)
	favoriteRepo := repository.NewFavoriteRepository(db)
	alertRepo := repository.NewAlertRepository(db)
	notifRepo := repository.NewNotificationRepository(db)
	postRepo := repository.NewPostRepository(db)
	followRepo := repository.NewFollowRepository(db)

	// Initialize services
	authSvc := service.NewAuthService(userRepo, jwt)
	articleSvc := service.NewArticleService(articleRepo, redisCache)
	exchangeSvc := service.NewExchangeRateService(exchangeRepo)
	likeSvc := service.NewLikeService(redisCache)
	rateHistorySvc := service.NewRateHistoryService(rateHistoryRepo, redisCache)
	favoriteSvc := service.NewFavoriteService(favoriteRepo)
	postSvc := service.NewPostService(postRepo, userRepo, followRepo)
	followSvc := service.NewFollowService(followRepo, userRepo)
	aiAnalystSvc := service.NewAIAnalystService(rateHistoryRepo, redisCache, cfg.LLM)

	// Initialize WebSocket hub
	hub := ws.NewHub()
	go hub.Run()

	// Initialize alert and notification services (depend on hub)
	alertSvc := service.NewAlertService(alertRepo, notifRepo, hub)
	notifSvc := service.NewNotificationService(notifRepo)

	// Initialize handlers
	h := router.Handlers{
		Auth:         handler.NewAuthHandler(authSvc),
		Article:      handler.NewArticleHandler(articleSvc),
		Exchange:     handler.NewExchangeHandler(exchangeSvc),
		Like:         handler.NewLikeHandler(likeSvc),
		RateHistory:  handler.NewRateHistoryHandler(rateHistorySvc),
		WS:           handler.NewWSHandler(hub),
		Favorite:     handler.NewFavoriteHandler(favoriteSvc),
		Alert:        handler.NewAlertHandler(alertSvc),
		Notification: handler.NewNotificationHandler(notifSvc),
		Post:         handler.NewPostHandler(postSvc),
		Follow:       handler.NewFollowHandler(followSvc),
		UserProfile:  handler.NewUserProfileHandler(userRepo),
		AIAnalyst:    handler.NewAIAnalystHandler(aiAnalystSvc),
	}

	// Setup router
	r, rateLimiter := router.SetupRouter(h, jwt, db, userRepo)

	// Start rate collector scheduler (with WebSocket broadcast)
	rateCollector := scheduler.NewRateCollector(rateHistorySvc, hub)
	rateCollector.Start()
	defer rateCollector.Stop()

	// Start server with graceful shutdown
	srv := &http.Server{
		Addr:         cfg.App.Port,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 120 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Printf("Server starting on %s", cfg.App.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")
	hub.Stop()
	rateLimiter.Stop()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited gracefully")
}
