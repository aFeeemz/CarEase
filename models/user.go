package models

type User struct {
	ID            uint    `gorm:"primaryKey"`
	Email         string  `gorm:"unique;not null" json:"email"`
	Username      string  `gorm:"unique;not null"`
	Password      string  `gorm:"not null"`
	DepositAmount float64 `gorm:"default:0" json:"deposit_amount"`
	IsAdmin       bool    `gorm:"column:isadmin;default:false"`
}

// To specify table name
func (User) TableName() string {
	return "users"
}

type LoginInfo struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

//
