package models

type User struct {
	ID            uint    `gorm:"primaryKey"`
	Email         string  `gorm:"unique;not null"`
	Username      string  `gorm:"unique;not null"`
	Password      string  `gorm:"not null"`
	DepositAmount float64 `gorm:"default:0"`
	IsAdmin       bool    `gorm:"default:false"`
}

type LoginInfo struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}
