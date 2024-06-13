package controllers

import (
	"FinalProject/initializers"
	"FinalProject/models"
	"net/http"

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

//
