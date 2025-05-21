package models

import "gorm.io/gorm"

type Product struct {
	gorm.Model
	Name        string
	Brand       string
	CategoryID  uint
	Price       float64
	Stock       uint
	Description string
}
