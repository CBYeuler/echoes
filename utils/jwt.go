package utils

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// GenerateJWT creates a new JWT token for a given username
// It uses the JWT_SECRET environment variable for signing the token
// Returns the signed token string or an error if something goes wrong
// The token is valid for 24 hours
// The JWT claims include the issuer, subject (username), and expiration time
// If the JWT_SECRET is not set, it returns an error
// If the token generation fails, it returns an error
// The token can be used for user authentication in the application

var ErrMissingSecret = errors.New("JWT_SECRET environment variable is not set")

func GenerateJWT(username string) (string, error) {
	//define the JWT claims
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return "", ErrMissingSecret
	}
	// Create the claims for the token
	claims := &jwt.RegisteredClaims{
		Issuer:    "echoes",
		Subject:   username,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	}
	//create the token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	//sign the token with the secret key
	signedToken, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}
	return signedToken, nil
}
func ValidateJWT(tokenString string) (*jwt.RegisteredClaims, error) {
	// Validate the JWT token and return the claims
	// It checks if the token is valid and returns the claims if successful
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return nil, ErrMissingSecret
	}
	// Parse the token and validate it
	// If the token is valid, it returns the claims
	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(secret), nil
	})
	// If there is an error parsing the token, return it
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*jwt.RegisteredClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("invalid JWT token")
}

//had troules with postman will try curl later

// This function can be used in the application to generate JWT tokens for users
