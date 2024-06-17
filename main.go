package main

import (
	"FinalProject/controllers"
	"FinalProject/initializers"
	"FinalProject/middleware"
	"fmt"

	"github.com/gin-gonic/gin"
)

func init() {
	initializers.LoadEnv()
	initializers.ConnectDB()
}

func main() {
	fmt.Println("hehe")

	r := gin.Default()

	// Public API
	customer := r.Group("/customer")
	{
		customer.POST("/register", controllers.Register)
		customer.POST("/login", controllers.Login)
		customer.GET("/carlist", middleware.AuthMiddlewareCustomer, controllers.GetAllAvailableCars)
		customer.POST("/rentcar", middleware.AuthMiddlewareCustomer, controllers.RentCar)
		customer.POST("/returncar", middleware.AuthMiddlewareCustomer, controllers.ReturnCar)

	}

	admin := r.Group("/admin")
	{
		admin.POST("/register", controllers.RegisterAdmin)
		admin.POST("/login", controllers.LoginAdmin)
		admin.POST("/addcar", middleware.AuthMiddleware, controllers.AddCars)
		admin.POST("/addcategory", middleware.AuthMiddleware, controllers.AddCategory)
	}
	//
	r.Run()

}
