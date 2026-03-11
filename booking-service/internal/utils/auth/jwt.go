package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/leminhthai/train-ticket/booking-service/global"
)

const (
	TokenTypeAccess  = "access"
	TokenTypeRefresh = "refresh"
)

type Claims struct {
	TokenType string `json:"token_type"`
	jwt.RegisteredClaims
}

func GenerateAccessToken(subjectUUID string) (string, error) {
	ttl := time.Duration(global.Config.JWT.ACCESS_TOKEN_TTL) * time.Minute
	return generateToken(subjectUUID, TokenTypeAccess, ttl)
}

func GenerateRefreshToken(subjectUUID string) (string, error) {
	ttl := time.Duration(global.Config.JWT.REFRESH_TOKEN_TTL) * time.Minute
	return generateToken(subjectUUID, TokenTypeRefresh, ttl)
}

func generateToken(subjectUUID, tokenType string, ttl time.Duration) (string, error) {
	secret := global.Config.JWT.API_SECRET
	if secret == "" {
		return "", fmt.Errorf("JWT secret is not configured")
	}

	now := time.Now()
	claims := Claims{
		TokenType: tokenType,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        uuid.New().String(),
			Subject:   subjectUUID,
			Issuer:    "train-ticket",
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(ttl)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func ParseToken(tokenStr string) (*Claims, error) {
	secret := global.Config.JWT.API_SECRET
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	if err != nil || !token.Valid {
		return nil, err
	}

	return claims, nil
}

func VerifyTokenSubject(tokenStr string) (*jwt.RegisteredClaims, error) {
	claims, err := ParseToken(tokenStr)
	if err != nil {
		return nil, err
	}

	// Chỉ cho phép Access Token đi qua middleware
	if claims.TokenType != TokenTypeAccess {
		return nil, fmt.Errorf("invalid token type")
	}

	return &claims.RegisteredClaims, nil
}