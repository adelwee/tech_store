package db

import (
	"Assignment2_AdelKenesova/user_service/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

var DB *gorm.DB

func InitDB() {
	dsn := "host=postgres user=postgres password=password dbname=tech_store port=5432 sslmode=disable"
	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to user_db: %v", err)
	}
	DB = database

	err = DB.AutoMigrate(&models.User{})
	if err != nil {
		log.Fatalf("Failed to migrate user table: %v", err)
	}

	log.Println("Connected and migrated user_db successfully")
}
