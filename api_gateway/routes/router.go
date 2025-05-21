package routes

import (
	"Assignment2_AdelKenesova/api_gateway/handlers"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	// Users
	r.POST("/users/register", handlers.RegisterUser)
	r.GET("/users/profile", handlers.GetUserProfile)

	// Products
	r.POST("/products", handlers.CreateProduct)
	r.GET("/products/:id", handlers.GetProduct)
	r.PATCH("/products/:id", handlers.UpdateProduct)
	r.DELETE("/products/:id", handlers.DeleteProduct)
	r.GET("/products", handlers.ListProducts)

	r.POST("/orders", handlers.CreateOrder)
	r.GET("/orders/:id", handlers.GetOrder)
	r.DELETE("/orders/:id", handlers.DeleteOrder)
	r.GET("/orders", handlers.ListOrders)

	return r
}
