package auth

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/leminhthai/train-ticket/user-service/global"
)

func ExtractBearerToken(c *gin.Context) (string, bool) {

	authHeader := c.GetHeader("Authorization")

	if strings.HasPrefix(authHeader, "Bearer ") {
		return strings.TrimPrefix(authHeader, "Bearer "), true
	}

	return "", false
}

func ParseJWTTokenSubject(token string) (*jwt.RegisteredClaims, error) {
	tokenClaims, err := jwt.ParseWithClaims(token, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(global.Config.JWT.API_SECRET), nil
	})

	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*jwt.RegisteredClaims); ok && tokenClaims.Valid {
			return claims, nil
		}
	}
	return nil, err
}
