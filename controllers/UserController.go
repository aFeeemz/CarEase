package controllers

import (
	"FinalProject/initializers"
	"FinalProject/models"
	"os"

	// "FinalProject/utils"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// Register handles user registration
func Register(c *gin.Context) {
	// Bind JSON input to the model
	var user models.User
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
	// Explicitly set IsAdmin to false
	user.IsAdmin = false

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
func Login(c *gin.Context) {
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
	var jwtKey = []byte(os.Getenv("SECRET_CUSTOMER"))

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour * 24).Unix(), //TOKEN ONLY WORKS FOR 5 MINUTES
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to Create Token",
		})

	}

	// Set cookies for username and token
	c.SetCookie("username", user.Username, int(time.Hour*24/time.Second), "/", "", false, true)
	c.SetCookie("token", tokenString, int(time.Hour*24/time.Second), "/", "", false, true)
	c.SetCookie("isAdmin", "false", int(time.Hour*24/time.Second), "/", "", false, true)

	// Respond with the generated token
	c.JSON(http.StatusOK, gin.H{
		"token": token,
	})
}
