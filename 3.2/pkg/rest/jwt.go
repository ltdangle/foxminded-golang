package rest

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt"
)

// Claims structure for JWT
type Claims struct {
	Role  string `json:"role"`
	Email string `json:"email"`
	jwt.StandardClaims
}

// JwtService interface
type JwtServiceInterface interface {
	GenerateToken(email string, role string) (string, error)
	ParseToken(tokenString string) (*Claims, error)
}

// JwtService for handling JWT operations
type JwtService struct {
	JwtKey []byte
}

func NewJwtService(key []byte) *JwtService {
	return &JwtService{JwtKey: key}
}

func (js *JwtService) GenerateToken(email string, role string) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		Role:  role,
		Email: email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(js.JwtKey)
}

func (js *JwtService) ParseToken(tokenString string) (*Claims, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return js.JwtKey, nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}
