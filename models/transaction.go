package models

import (
	"time"
)

type Transaction struct {
	ID            uint    `gorm:"primaryKey"`
	Email         string  `gorm:"unique;not null" json:"email"`
	Username      string  `gorm:"unique;not null"`
	Password      string  `gorm:"not null"`
	DepositAmount float64 `gorm:"default:0" json:"deposit_amount"`
}

type RentalRequest struct {
	CarID      uint      `json:"car_id" binding:"required"`
	ReturnDate time.Time `json:"return_date" binding:"required"`
}

type RentalHistory struct {
	UserID           uint      `gorm:"not null"`
	CarID            uint      `gorm:"not null"`
	RentalDate       time.Time `gorm:"not null"`
	ReturnDate       time.Time `gorm:"not null"`
	ActualReturnDate time.Time
	TotalCost        float64 `gorm:"not null"`
	PenaltyRate      float64 `gorm:"not null;default:10.0"` // Example penalty rate per day
	PenaltyAmount    float64 `gorm:"default:0"`
	Paid             bool    `gorm:"not null;default:false"`
	User             User
	Car              Car
}

func (RentalHistory) TableName() string {
	return "rentalhistory"
}

type ReturnRequest struct {
	RentalHistoryID uint `json:"rental_history_id" binding:"required"`
}

//
