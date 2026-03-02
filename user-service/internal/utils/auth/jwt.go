package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/leminhthai/train-ticket/user-service/global"
)

type Claims struct {
	jwt.RegisteredClaims
}

func GenerateToken(uuidToken string) (string, error) {
	secret := global.Config.JWT.API_SECRET
	timeEx := global.Config.JWT.JWT_EXPIRATION

	expiration, err := time.ParseDuration(timeEx)

	if err != nil {
		return "", fmt.Errorf("Invalid JWT_EXPIRATION duration: %s, error: %w", timeEx, err)
	}

	now := time.Now()
	expiresAt := now.Add(expiration)

	claims := Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        uuid.New().String(),
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(now),
			Subject:   uuidToken,
			Issuer:    "train-ticket",
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

func VerifyTokenSubject(token string) (*jwt.RegisteredClaims, error) {
	claims, err := ParseJWTTokenSubject(token)
	if err != nil {
		return nil, err
	}
	return claims, nil
}
