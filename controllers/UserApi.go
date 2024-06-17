package controllers

import (
	"FinalProject/initializers"
	"FinalProject/models"
	"FinalProject/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func GetAllAvailableCars(c *gin.Context) {
	var cars []models.Car
	// Query available cars
	err := initializers.DB.Where("availability = ?", true).Find(&cars).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch available cars"})
		return
	}

	// Return JSON response
	c.JSON(http.StatusOK, gin.H{"cars": cars})
}

type RentCarInput struct {
	CarID      uint      `json:"car_id" binding:"required"`
	RentalDate time.Time `json:"rental_date" binding:"required"`
	ReturnDate time.Time `json:"return_date" binding:"required"`
}

// RentCar handles the car rental process
func RentCar(c *gin.Context) {
	var input RentCarInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, exists := utils.GetUserFromContext(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	db := initializers.DB

	// Check if the car is available
	var car models.Car
	if err := db.First(&car, input.CarID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Car not found"})
		return
	}
	if !car.Availability {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Car is not available"})
		return
	}

	// Calculate the total cost based on the rental period
	rentalDuration := input.ReturnDate.Sub(input.RentalDate).Hours() / 24
	totalCost := rentalDuration * car.RentalCosts

	// Create the rental history record
	rentalHistory := models.RentalHistory{
		UserID:     user.ID,
		CarID:      input.CarID,
		RentalDate: input.RentalDate,
		ReturnDate: input.ReturnDate,
		TotalCost:  totalCost,
	}

	if err := utils.CreateRentalHistory(rentalHistory); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create rental history"})
		return
	}

	// Update the car availability
	car.Availability = false
	if err := db.Save(&car).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update car availability"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": rentalHistory})
}

func ReturnCar(c *gin.Context) {
	var input struct {
		RentalHistoryID  uint      `json:"rental_history_id" binding:"required"`
		ActualReturnDate time.Time `json:"actual_return_date" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON input"})
		return
	}

	db := initializers.DB

	// Retrieve the rental history record using raw SQL query
	var rentalHistory models.RentalHistory
	queryRentalHistory := `
	SELECT id, user_id, car_id, rental_date, return_date, actual_return_date, total_cost, penalty_amount, paid
	FROM rentalhistory
	WHERE id = $1
`
	if err := db.Raw(queryRentalHistory, input.RentalHistoryID).Scan(&rentalHistory).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Rental history not found"})
		return
	}

	// Check if the rental history has already been paid
	if rentalHistory.Paid {
		c.JSON(http.StatusBadRequest, gin.H{"error": "The bill has already been paid"})
		return
	}

	// Check if the actual return date is after the expected return date
	if input.ActualReturnDate.After(rentalHistory.ReturnDate) {
		// Calculate late days
		lateDays := input.ActualReturnDate.Sub(rentalHistory.ReturnDate).Hours() / 24
		// Calculate penalty (adjust penalty calculation as per your business logic)
		penalty := lateDays * 50 // Assuming $50 penalty per late day

		// Update penalty amount in rental history
		rentalHistory.PenaltyAmount = penalty
		// Update total cost with penalty
		rentalHistory.TotalCost += penalty
	}

	// Update the actual return date and penalty amount in rental history
	rentalHistory.ActualReturnDate = input.ActualReturnDate

	// Deduct total cost from user's deposit_amount
	user, exists := utils.GetUserFromContext(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	if user.DepositAmount < rentalHistory.TotalCost {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Insufficient deposit amount"})
		return
	}

	user.DepositAmount -= rentalHistory.TotalCost
	// Begin a transaction
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Execute raw SQL query to update rental history
	query := `
		UPDATE rentalhistory
		SET actual_return_date = $1, penalty_amount = $2, total_cost = $3, paid = true
		WHERE id = $4
	`
	if err := tx.Exec(query, rentalHistory.ActualReturnDate, rentalHistory.PenaltyAmount, rentalHistory.TotalCost, input.RentalHistoryID).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update rental history"})
		return
	}

	// Update user's deposit amount
	if err := tx.Model(&user).Update("deposit_amount", user.DepositAmount).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user deposit amount"})
		return
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Car returned successfully"})
}
