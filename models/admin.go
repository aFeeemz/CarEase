package models

type Admin struct {
	ID            uint    `gorm:"primaryKey"`
	Email         string  `gorm:"unique;not null" json:"email"`
	Username      string  `gorm:"unique;not null"`
	Password      string  `gorm:"not null"`
	DepositAmount float64 `gorm:"default:0" json:"deposit_amount"`
	IsAdmin       bool    `gorm:"column:isadmin;default:false"`
}

// To specify table name
func (Admin) TableName() string {
	return "users"
}

type Car struct {
	ID           uint    `gorm:"primaryKey"`
	Availability bool    `gorm:"not null;default:true"`
	RentalCosts  float64 `gorm:"not null"`
	CategoryID   uint    `gorm:"not null"`
}

// To specify table name
func (Car) TableName() string {
	return "cars"
}

type Category struct {
	ID           uint   `gorm:"primaryKey"`
	Name         string `gorm:"not null"`
	VehicleBrand string `gorm:"not null" json:"vehicle_brand"`
	Color        string `gorm:"not null"`
	Transmission string `gorm:"not null"`
	VinNumber    string `gorm:"not null" json:"vin_number"`
}

// To specify table name
func (Category) TableName() string {
	return "category"
}

//
