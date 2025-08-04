package main

import (
	"fmt"
	"log"
	"os"

	"github.com/CBYeuler/echoes/config"
	"github.com/CBYeuler/echoes/controllers"
	"github.com/CBYeuler/echoes/database"
	"github.com/CBYeuler/echoes/middleware"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system env variables")
	}
	// Ensure the OpenAI API key is set
	// This key is required for making requests to OpenAI's API
	apiKey := os.Getenv("OPENAI_API_KEY")

	if apiKey == "" {
		log.Fatal("OPENAI_API_KEY is not set")
	}

	port := config.GetPort()
	if port <= 0 {
		log.Fatal("Invalid port number")
	}
	// Initialize the database
	database.InitDB()
	// Init router
	router := gin.Default()
	// Set up routes
	router.GET("/", func(c *gin.Context) {
		c.String(200, "Welcome to the Echoes API!")
	})
	// Health check endpoint
	// This endpoint can be used to check if the API is running
	// It returns a simple message indicating the API is healthy
	router.GET("/health", func(c *gin.Context) {
		c.String(200, "API is healthy")
	})
	// API key endpoint
	// This endpoint returns the API key used for OpenAI requests
	router.GET("/api-key", func(c *gin.Context) {
		c.String(200, fmt.Sprintf("Your API Key is: %s", apiKey))
	})

	//---------------------------------------------------------------------
	// Authentication routes
	// These routes handle user registration and login

	// grouping them under /auth for better organization
	auth := router.Group("/auth")
	{
		auth.POST("/register", controllers.Register)
		auth.POST("/login", controllers.Login)
	}
	//---------------------------------------------------------------------
	// Protected routes
	// These routes require authentication via JWT
	api := router.Group("/api")
	api.Use(middleware.AuthMiddleware()) // JWT-protected routes
	{
		api.POST("/echo", controllers.HandleEcho)
	}
	//add routes for user management, message history, etc.

	// Start the server

	log.Printf(" Echoes running on port %d", port)
	if err := router.Run(fmt.Sprintf(":%d", port)); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
