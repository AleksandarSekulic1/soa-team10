// encounters-service/api/middleware.go
package api

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var jwtSecretKey = []byte("super_secret_key") // Pretpostavka da ste uneli pravi ključ

type AppClaims struct {
	Id       int64  `json:"id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is missing"})
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Bearer token not found"})
			return
		}

		claims := &AppClaims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return jwtSecretKey, nil
		})

		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			return
		}

		// --- ISPRAVKA OVDE ---
		// Postavljamo 'username' kao string za 'userId' da bi tipovi bili usklađeni
		c.Set("userId", claims.Username)
		
		// Ostavljamo i 'username' i 'role' za svaki slučaj, ako zatrebaju
		c.Set("username", claims.Username)
		c.Set("role", claims.Role)

		c.Next()
	}
}