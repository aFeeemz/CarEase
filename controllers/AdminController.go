package controllers

import (
	"FinalProject/initializers"
	"FinalProject/models"
	"FinalProject/utils"
	"fmt"
	"net/http"
	"time"

	// "time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// Register handles user registration
func RegisterAdmin(c *gin.Context) {
	// Bind JSON input to the model
	var user models.Admin
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON input"})
		return
	}

	// Hashing the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}
	user.Password = string(hashedPassword)

	// Explicitly declare isadmin to true
	user.IsAdmin = true

	// Explicitly declare wallet to zero
	user.DepositAmount = 0

	// Make the user to the database
	if err := initializers.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Respond with success message
	c.JSON(http.StatusOK, gin.H{
		"success": fmt.Sprintf("User %v has been successfully registered", user.Username),
	})
}

// Login handles user login
func LoginAdmin(c *gin.Context) {
	fmt.Println("MASUK LOGIN ADMIN")
	// Bind JSON input to the model
	var loginInfo models.LoginInfo
	if err := c.ShouldBindJSON(&loginInfo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON input"})
		return
	}

	// Find the user by email
	var user models.User
	result := initializers.DB.Where("username = ?", loginInfo.Username).First(&user)
	if result.Error != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Username"})
		return
	}

	// Check if the password is correct
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginInfo.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Password"})
		return
	}

	// Generate JWT token for the user
	token, err := utils.GenerateJWT(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})

	}

	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("username", user.Username, int(time.Hour*24/time.Second), "/", "", false, true)
	c.SetCookie("token", token, int(time.Hour*24/time.Second), "/", "", false, true)
	c.SetCookie("isAdmin", "true", int(time.Hour*24/time.Second), "/", "", false, true)

	if user.IsAdmin {
		c.SetCookie("isAdmin", "true", int(time.Hour*24/time.Second), "", "", false, true)
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user should be admin"})
	}

	// Respond with the generated token
	c.JSON(http.StatusOK, gin.H{
		"token": token,
	})
}
