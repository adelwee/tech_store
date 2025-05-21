package models

import (
	"gorm.io/gorm"
)

type Order struct {
	gorm.Model
	UserID     uint
	TotalPrice float64
	Status     string
	OrderItems []OrderItem `gorm:"foreignKey:OrderID"`
}

type OrderItem struct {
	gorm.Model
	OrderID   uint
	ProductID uint
	Quantity  uint
	Price     float64
}
