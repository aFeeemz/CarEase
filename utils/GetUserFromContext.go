package utils

// // GetUserFromContext retrieves the user from the Gin context
// func GetUserFromContext(c *gin.Context) (*models.User, error) {
// 	user, exists := c.Get("user")
// 	if !exists {
// 		return nil, errors.New("user not found in context")
// 	}

// 	userModel, ok := user.(*models.User)
// 	if !ok {
// 		return nil, errors.New("user is of invalid type")
// 	}

// 	return userModel, nil
// }
