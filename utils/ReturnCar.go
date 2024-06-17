package utils

import (
	"FinalProject/initializers"
	"FinalProject/models"
	"time"
)

type CreateRentalHistoryInput struct {
	UserID     uint      `json:"user_id" binding:"required"`
	CarID      uint      `json:"car_id" binding:"required"`
	RentalDate time.Time `json:"rental_date" binding:"required"`
	ReturnDate time.Time `json:"return_date" binding:"required"`
	TotalCost  float64   `json:"total_cost" binding:"required"`
}

// CreateRentalHistory creates a new rental history record
func CreateRentalHistory(rentalHistory models.RentalHistory) error {
	db := initializers.DB
	if err := db.Create(&rentalHistory).Error; err != nil {
		return err
	}
	return nil
}
