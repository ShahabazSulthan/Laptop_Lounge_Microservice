package server

import (
	"api-gateway/pkg/handler"
	"api-gateway/pkg/middleware"
	"log"

	"github.com/gin-gonic/gin"
)

type ServerHTTP struct {
	engine *gin.Engine
}

// NewServerHTTP initializes the server with routes and handlers
func NewServerHTTP(adminHandler *handler.AdminHandler, productHandler *handler.ProductHandler, userHandler *handler.UserHandler, cartHandler *handler.CartHandler, orderHandler *handler.OrderHandler) *ServerHTTP {
	router := gin.New()

	router.Use(gin.Logger())
	router.Use(gin.Recovery()) // Add recovery middleware to handle panics

	// Public routes
	router.POST("/admin/login", adminHandler.LoginHandler)
	router.POST("/admin/signup", adminHandler.AdminSignUp)
	router.POST("/user/signup", userHandler.UserSignup)
	router.POST("/user/login", userHandler.Userlogin)
	router.GET("/product", productHandler.ShowAllProducts)

	// Admin routes
	adminRoutes := router.Group("/")
	adminRoutes.Use(middleware.AdminAuthMiddleware())
	{
		adminRoutes.POST("/product", productHandler.AddProducts)
		adminRoutes.DELETE("/product", productHandler.DeleteProduct)
		adminRoutes.PUT("/product", productHandler.UpdateProducts)
	}

	// User routes
	userRoutes := router.Group("/")
	userRoutes.Use(middleware.UserAuthMiddleware())
	{
		userRoutes.POST("/cart", cartHandler.AddToCart)
		userRoutes.GET("/cart", cartHandler.GetCart)
		userRoutes.POST("/order", orderHandler.OrderItemsFromCart)
		userRoutes.GET("/order", orderHandler.GetOrderDetails)

		// Address routes
		userRoutes.POST("/address", userHandler.AddAddress)
		userRoutes.GET("/address/:id", userHandler.GetAddress)
		userRoutes.PUT("/address", userHandler.UpdateAddress)
		userRoutes.DELETE("/address/:id", userHandler.DeleteAddress)
	}

	return &ServerHTTP{engine: router}
}

// Start runs the server on the specified port
func (s *ServerHTTP) Start() {
	log.Println("Starting server on :6000")
	err := s.engine.Run(":6000")
	if err != nil {
		log.Fatalf("Error while starting the server: %v", err)
	}
}
