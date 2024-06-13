package controllers

import (
	"FinalProject/initializers"
	"FinalProject/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// AddCars adds a new car to the database
func AddCars(c *gin.Context) {
	// Check if the user is an admin
	isAdmin, exists := c.Get("isAdmin")
	if !exists || !isAdmin.(bool) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only admins can add cars"})
		return
	}

	// Bind JSON input to the model
	var input struct {
		RentalCosts float64 `json:"rental_costs" binding:"required"`
		CategoryID  uint    `json:"category_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON input"})
		return
	}

	// Create a new car entry
	car := models.Car{
		Availability: true,
		RentalCosts:  input.RentalCosts,
		CategoryID:   input.CategoryID,
	}

	// Save the car to the database
	if err := initializers.DB.Create(&car).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Respond with success message
	c.JSON(http.StatusOK, gin.H{"success": "Car has been successfully added"})
}

// AddCategory handles the addition of new categories by admin
func AddCategory(c *gin.Context) {
	// Check if the user is an admin
	isAdmin, exists := c.Get("isAdmin")
	if !exists || !isAdmin.(bool) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only admins can add categories"})
		return
	}

	// Bind JSON input to the model
	var category models.Category
	if err := c.ShouldBindJSON(&category); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON input"})
		return
	}

	// Save the category to the database
	if err := initializers.DB.Create(&category).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Respond with success message
	c.JSON(http.StatusOK, gin.H{"success": "Category has been successfully added"})
}

//
