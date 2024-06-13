package middleware

import (
	"FinalProject/initializers"
	"FinalProject/models"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// AuthMiddleware checks the JWT token and verifies if the user is an admin
func AuthMiddleware(c *gin.Context) {
	tokenString, err := c.Cookie("token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized, no token provided"})
		c.Abort()
		return
	}

	// Decode and validate
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("SECRET_CUSTOMER")), nil
	})
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		c.Abort()
		return
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Check expiration
		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token expired"})
			c.Abort()
			return
		}

		// Check isAdmin claim
		if isAdmin, ok := claims["isAdmin"].(bool); ok && isAdmin {
			// Set user ID and isAdmin to the context
			c.Set("userID", claims["sub"])
			c.Set("isAdmin", true)
		} else {
			c.JSON(http.StatusForbidden, gin.H{"error": "Only admins can access this resource"})
			c.Abort()
			return
		}

		// Find user with token
		var user models.User
		initializers.DB.First(&user, claims["sub"])
		if user.ID == 0 {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
			c.Abort()
			return
		}

		// Attach user to the context
		c.Set("user", user)

		// Continue to the next handler
		c.Next()
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
		c.Abort()
		return
	}
}

// func AuthMiddlewareAdmin(c *gin.Context) {

// 	// Get the token from the cookie
// 	token, err := c.Cookie("token")
// 	if err != nil {
// 		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized, no token provided"})
// 		c.Abort()
// 		return
// 	}

// 	// Validate the token
// 	claims, err := utils.ValidateJWTadmin(token)
// 	if err != nil {
// 		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized, invalid admin token "})
// 		c.Abort()
// 		return
// 	}

// 	// Set user ID to the context
// 	c.Set("userID", claims.UserID)

// 	fmt.Println("PASSED MIDDLEWARE AUTH ADMIN")
// 	// Continue to the next handler
// 	c.Next()

// }
