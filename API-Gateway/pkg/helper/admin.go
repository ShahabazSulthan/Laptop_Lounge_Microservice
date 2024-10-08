package helper

import (
	"api-gateway/pkg/utils/models"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
)

// Key for signing the JWT. Consider loading this from an environment variable or a config file.
var jwtKey = []byte("admin_token")

// authCustomClaims represents the JWT claims.
type authCustomClaims struct {
	Firstname string `json:"firstname,omitempty"`
	Lastname  string `json:"lastname,omitempty"`
	Email     string `json:"email"`
	Role      string `json:"role"` // Added role field
	jwt.StandardClaims
}

// GenerateToken generates a JWT token with role information.
func GenerateToken(admin models.AdminDetailsResponse) (string, error) {
	claims := &authCustomClaims{
		Firstname: admin.Firstname,
		Lastname:  admin.Lastname,
		Email:     admin.Email,
		Role:      "admin", // Assign role here
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(48 * time.Hour).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		fmt.Println("Error signing the token:", err)
		return "", err
	}

	return tokenString, nil
}

// ValidateToken validates a JWT token and returns the claims if valid.
func ValidateToken(tokenString string) (*authCustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &authCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtKey, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*authCustomClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, fmt.Errorf("invalid token")
}
