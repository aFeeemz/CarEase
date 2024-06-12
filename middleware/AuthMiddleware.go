package middleware

import (
	"FinalProject/utils"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the token from the cookie
		token, err := c.Cookie("token")
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized, no token provided"})
			c.Abort()
			return
		}

		// Validate the token
		claims, err := utils.ValidateJWT(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized, invalid token"})
			c.Abort()
			return
		}

		// Set user ID to the context
		c.Set("userID", claims.UserID)

		// Continue to the next handler
		c.Next()
	}
}

func AuthMiddlewareAdmin(c *gin.Context) {

	// Get the token from the cookie
	token, err := c.Cookie("token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized, no token provided"})
		c.Abort()
		return
	}

	// Validate the token
	claims, err := utils.ValidateJWT(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized, invalid token"})
		c.Abort()
		return
	}

	// Set user ID to the context
	c.Set("userID", claims.UserID)

	fmt.Println("PASSED MIDDLEWARE AUTH ADMIN")
	// Continue to the next handler
	c.Next()

}
