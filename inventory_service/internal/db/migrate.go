package db

import (
	"log"
)

func Migrate() {
	err := DB.AutoMigrate(&Product{})
	if err != nil {
		log.Fatal("Failed to run migrations:", err)
	}
	log.Println("Migrations completed successfully!")
}
