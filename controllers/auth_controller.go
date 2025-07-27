package controllers

import (
	"log"
	"net/http"

	"github.com/CBYeuler/echoes/database"
	"github.com/CBYeuler/echoes/models"
	"github.com/CBYeuler/echoes/utils"
	"github.com/gin-gonic/gin"
)

func Register(c *gin.Context) {
	// Register handles user registration
	// It expects a JSON body with username and password
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

func Login(c *gin.Context) {
	// Login handles user login
	// It checks the username and password, generates a JWT token if successful

	var input models.User
	if err := c.ShouldBindJSON(&input); err != nil {
		log.Println("Failed to bind JSON input:", err) // Log binding errors
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	//lookup the user in the database
	query := "SELECT id, password FROM users WHERE username = ?"
	row := database.DB.QueryRow(query, input.Username)
	var user models.User
	if err := row.Scan(&user.ID, &user.Password); err != nil {
		log.Println("User not found or error scanning row:", err) // Log if user not found or scan error
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	// Check if the provided password matches the stored hash
	if !models.CheckPasswordHash(input.Password, user.Password) {
		log.Println("Password mismatch for user:", input.Username) // Log password mismatch
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	// Generate JWT token for the user
	token, err := utils.GenerateJWT(input.Username)
	if err != nil {
		log.Println("Failed to generate JWT token:", err) // Log JWT generation errors
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}
	// Log the successful login and return the token
	log.Println("User logged in successfully:", input.Username) // Log successful login
	c.JSON(http.StatusOK, gin.H{"token": token})                // Return the JWT token
	return
}
