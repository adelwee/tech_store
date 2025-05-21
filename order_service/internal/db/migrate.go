package db

import (
	"Assignment2_AdelKenesova/order_service/internal/models"
	"log"
)

func Migrate() {
	err := DB.AutoMigrate(&models.Order{}, &models.OrderItem{})
	if err != nil {
		log.Fatalf("Failed to migrate OrderService tables: %v", err)
	}
	log.Println("OrderService tables migrated successfully")
}
