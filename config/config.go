package config

import (
	"os"
	"strconv"
)

func GetPort() int {
	// Read the PORT environment variable
	// If not set, return a default port
	port := os.Getenv("PORT")
	if port == "" {
		return 8080 // Default port
	}

	p, err := strconv.Atoi(port)
	if err != nil {
		return 8080 // Fallback to default if conversion fails
	}
	return p
}
