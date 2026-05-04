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
	); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	log.Println("Migrations completed successfully")
}
