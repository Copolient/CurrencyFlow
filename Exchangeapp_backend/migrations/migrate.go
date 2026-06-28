package migrations

import (
	"exchangeapp/internal/model"
	"log"

	"gorm.io/gorm"
)

// Run executes AutoMigrate for all models.
// WARNING: AutoMigrate can add columns and indexes but cannot drop or rename them.
// For production, consider using a versioned migration tool (e.g., golang-migrate).
func Run(db *gorm.DB) {
	log.Println("Running database migrations...")

	if err := db.AutoMigrate(
		&model.User{},
		&model.Article{},
		&model.ExchangeRate{},
		&model.ExchangeRateHistory{},
		&model.Favorite{},
		&model.RateAlert{},
		&model.Notification{},
		&model.Post{},
		&model.PostLike{},
		&model.Follow{},
	); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	log.Println("Migrations completed successfully")
}
