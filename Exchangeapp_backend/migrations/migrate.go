package migrations

import (
	"exchangeapp/internal/model"
	"log"

	"gorm.io/gorm"
)

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
		&model.Follow{},
	); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	log.Println("Migrations completed successfully")
}
