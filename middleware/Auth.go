package middleware

import (
	"log"
	"net/http"
	"strings"

	"github.com/CBYeuler/echoes/utils"
	"github.com/gin-gonic/gin"
)

// AuthMiddleware is a middleware function that checks for a valid JWT token in the request header

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			log.Println("Authorization header is missing")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			c.Abort()
			return
		}

		// Check if the header starts with "Bearer "
		if !strings.HasPrefix(authHeader, "Bearer ") {
			log.Println("Authorization header format is invalid")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization format"})
			c.Abort()
			return
		}

		// Extract the token from the header
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// Validate the JWT token
		claims, err := utils.ValidateJWT(tokenString)
		if err != nil {
			log.Println("Invalid JWT token:", err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// Store the claims in the context for use in handlers
		c.Set("claims", claims)
		c.Set("username", claims.Subject)
		c.Next()
	}
}
