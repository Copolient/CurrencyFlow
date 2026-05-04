package main

import (
	"context"
	"exchangeapp/internal/handler"
	"exchangeapp/internal/repository"
	"exchangeapp/internal/service"
	"exchangeapp/migrations"
	"exchangeapp/pkg/auth"
	"exchangeapp/pkg/cache"
	"exchangeapp/pkg/config"
	"exchangeapp/pkg/database"
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

	// Initialize services
	authSvc := service.NewAuthService(userRepo, jwt)
	articleSvc := service.NewArticleService(articleRepo, redisCache)
	exchangeSvc := service.NewExchangeRateService(exchangeRepo)
	likeSvc := service.NewLikeService(redisCache)

	// Initialize handlers
	h := router.Handlers{
		Auth:     handler.NewAuthHandler(authSvc),
		Article:  handler.NewArticleHandler(articleSvc),
		Exchange: handler.NewExchangeHandler(exchangeSvc),
		Like:     handler.NewLikeHandler(likeSvc),
	}

	// Setup router
	r := router.SetupRouter(h, jwt, db)

	// Start server with graceful shutdown
	srv := &http.Server{
		Addr:    cfg.App.Port,
		Handler: r,
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

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited gracefully")
}
