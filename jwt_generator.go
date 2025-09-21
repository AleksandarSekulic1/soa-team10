package main

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtKey = []byte("super_secret_key")

type Claims struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

func main() {
	// Kreiraj claims za test korisnika
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		ID:       "user1",
		Username: "testuser1",
		Role:     "tourist",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		fmt.Printf("Error creating token: %v\n", err)
		return
	}

	fmt.Printf("JWT Token for user1: %s\n", tokenString)
	
	// Kreiraj i drugi token za user2
	claims2 := &Claims{
		ID:       "user2",
		Username: "testuser2",
		Role:     "tourist",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token2 := jwt.NewWithClaims(jwt.SigningMethodHS256, claims2)
	tokenString2, err := token2.SignedString(jwtKey)
	if err != nil {
		fmt.Printf("Error creating token: %v\n", err)
		return
	}

	fmt.Printf("JWT Token for user2: %s\n", tokenString2)
}