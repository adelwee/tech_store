package db

type Product struct {
	ID          uint64 `gorm:"primaryKey"`
	Name        string
	Brand       string
	CategoryID  uint64
	Price       float64
	Stock       uint64
	Description string
}
