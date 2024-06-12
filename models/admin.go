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
