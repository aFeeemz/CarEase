package controllers

import (
	"FinalProject/initializers"
	"FinalProject/models"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	// "time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
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

// LoginAdmin handles admin login
func LoginAdmin(c *gin.Context) {
	// Bind JSON input to the model
	var loginInfo models.LoginInfo
	if err := c.ShouldBindJSON(&loginInfo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON input"})
		return
	}

	// Find the user by username
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

	// Determine isAdmin status (assuming IsAdmin is a boolean field in your User model)
	isAdmin := user.IsAdmin

	// Generate JWT token for the user
	var jwtKey = []byte(os.Getenv("SECRET_CUSTOMER"))

	fmt.Printf("INI IS ADMIN DI ADMINCONTROLLER : %v\n", isAdmin)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":     user.ID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
		"isAdmin": isAdmin,
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to Create Token"})
		return
	}

	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("username", user.Username, int(time.Hour*24/time.Second), "/", "", false, true)
	c.SetCookie("token", tokenString, int(time.Hour*24/time.Second), "/", "", false, true)
	c.SetCookie("isAdmin", strconv.FormatBool(isAdmin), int(time.Hour*24/time.Second), "/", "", false, true)

	// Respond with the generated token
	c.JSON(http.StatusOK, gin.H{"token": tokenString})
}

//
