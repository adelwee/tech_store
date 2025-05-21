package db

import (
	"Assignment2_AdelKenesova/order_service/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

var DB *gorm.DB

func InitDB() {
	dsn := "host=localhost user=postgres password=0000 dbname=order_db port=5432 sslmode=disable"
	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to Order DB:", err)
	}
	DB = database

	err = DB.AutoMigrate(&models.Order{}, &models.OrderItem{})
	if err != nil {
		log.Fatal("Failed to migrate Order tables:", err)
	}

	log.Println("Connected and migrated Order DB successfully.")
}
