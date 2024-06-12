package models

type Car struct {
	ID           uint    `gorm:"primaryKey"`
	Name         string  `gorm:"not null"`
	Availability bool    `gorm:"not null;default:true"`
	RentalCosts  float64 `gorm:"not null"`
	VINNumber    string  `gorm:"not null;unique"`
	Color        string  `gorm:"not null"`
	Transmission string  `gorm:"not null"`
	CategoryID   uint    `gorm:"not null"`
}
