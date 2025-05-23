package db

import (
	"Assignment2_AdelKenesova/inventory_service/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

var DB *gorm.DB

func InitDB() {
	dsn := "host=postgres user=postgres password=password dbname=tech_store port=5432 sslmode=disable"
	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	DB = database
	log.Println("Connected to PostgreSQL successfully")

	err = database.AutoMigrate(&models.Product{})
	if err != nil {
		log.Fatal("Failed to migrate database: ", err)
	}
	log.Println("Migrated Product table successfully")

}
