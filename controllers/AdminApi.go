package controllers

import (
	"FinalProject/initializers"
	"FinalProject/models"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AddCars(c *gin.Context) {
	isAdmin, exists := c.Get("isAdmin")
	if !exists || !isAdmin.(bool) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access forbidden. User is not an admin"})
		return
	}

	// Bind JSON input to the model
	var car models.Car
	if err := c.ShouldBindJSON(&car); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON input"})
		return
	}

	// Create the car in the database
	if err := initializers.DB.Create(&car).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Respond with success message
	c.JSON(http.StatusOK, gin.H{
		"success": fmt.Sprintf("Car %v has been successfully added", car.Name),
	})
}
