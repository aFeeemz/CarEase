package utils

import (
	"FinalProject/initializers"
	"FinalProject/models"

	"github.com/gin-gonic/gin"
)

// getCarByID retrieves a car from the database by its ID
func GetCarByID(carID uint) (*models.Car, error) {
	var car models.Car
	if err := initializers.DB.Preload("Category").First(&car, carID).Error; err != nil {
		return nil, err
	}
	return &car, nil
}

func GetUserFromContext(c *gin.Context) (models.User, bool) {
	user, exists := c.Get("user")
	if !exists {
		return models.User{}, false
	}
	return user.(models.User), true
}

// getRentalHistoryByID retrieves a rental history record by its ID
func GetRentalHistoryByID(rentalHistoryID uint) (models.RentalHistory, error) {
	var rentalHistory models.RentalHistory
	db := initializers.DB
	if err := db.Preload("Car").Preload("User").First(&rentalHistory, rentalHistoryID).Error; err != nil {
		return models.RentalHistory{}, err
	}
	return rentalHistory, nil
}
