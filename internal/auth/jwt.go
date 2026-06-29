package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	defaultTokenDuration = 24 * time.Hour
)

type JWTClaims struct {
	UserID uint   `json:"user_id"`
	Name   string `json:"name"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

type JWTService interface {
	GenerateToken(userId uint, email, name, role string) (string, error)
	ValidateToken(tokenStr string) (*JWTClaims, error)
}

type jwtService struct {
	secretKey     string
	tokenDuration time.Duration
}

func NewJWTService(secretKey string) JWTService {
	return &jwtService{
		secretKey:     secretKey,
		tokenDuration: defaultTokenDuration,
	}
}

func (js *jwtService) GenerateToken(userId uint, email, name, role string) (string, error) {
	claims := JWTClaims{
		UserID: userId,
		Name:   name,
		Email:  email,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(js.tokenDuration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "spotsync",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(js.secretKey))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func (js *jwtService) ValidateToken(tokenStr string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &JWTClaims{}, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(js.secretKey), nil
	})

	if err != nil {
		return nil, fmt.Errorf("invalid or expired token: %w", err)
	}

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}
