package main

import (
	"fmt"
	"log"
	"os"

	"github.com/CBYeuler/echoes/config"
	"github.com/CBYeuler/echoes/controllers"
	"github.com/CBYeuler/echoes/database"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system env variables")
	}

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

	router.GET("/", func(c *gin.Context) {
		c.String(200, "Welcome to the Echoes API!")
	})

	router.GET("/health", func(c *gin.Context) {
		c.String(200, "API is healthy")
	})

	router.GET("/api-key", func(c *gin.Context) {
		c.String(200, fmt.Sprintf("Your API Key is: %s", apiKey))
	})
	router.POST("/register", controllers.Register)

	log.Printf(" Echoes running on port %d", port)
	if err := router.Run(fmt.Sprintf(":%d", port)); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
