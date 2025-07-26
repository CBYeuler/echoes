package controllers

import (
	"log"
	"net/http"

	"github.com/CBYeuler/echoes/database"
	"github.com/CBYeuler/echoes/models"
	"github.com/gin-gonic/gin"
)

// Register handles user registration
func Register(c *gin.Context) {
	var user models.User

	// Attempt to bind JSON input to the user struct
	if err := c.ShouldBindJSON(&user); err != nil {
		log.Println("Failed to bind JSON input:", err) //  Log what's wrong with the request body
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// Hash the password before storing it
	hashedPassword, err := models.HashPassword(user.Password)
	if err != nil {
		log.Println("Error hashing password:", err) //  Log hashing failures (e.g., bcrypt issues)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}
	user.Password = hashedPassword

	// Insert the user into the database
	query := "INSERT INTO users (username, password) VALUES (?, ?)"
	_, err = database.DB.Exec(query, user.Username, user.Password)
	if err != nil {
		log.Println("Failed to insert user into DB:", err) // Log DB errors (e.g., duplicate usernames)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register user"})
		return
	}

	log.Println("User registered successfully:", user.Username) // Log successful registration
	c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully"})
}
