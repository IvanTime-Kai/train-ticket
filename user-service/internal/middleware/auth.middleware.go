package middleware

import (
	"context"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/leminhthai/train-ticket/user-service/internal/utils/auth"
	"github.com/leminhthai/train-ticket/user-service/internal/utils/cache"
	"github.com/leminhthai/train-ticket/user-service/pkg/response"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		uri := c.Request.URL.Path
		log.Println("URI request at auth middleware", uri)

		jwtToken, ok := auth.ExtractBearerToken(c)

		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code": response.ErrUnauthorized,
				"err":  "Unauthorized", "description": "",
			})
			return
		}

		claims, err := auth.VerifyTokenSubject(jwtToken)

		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code": response.ErrUnauthorized,
				"err":  "Invalid token", "description": "",
			})
			return
		}

		// check black list
		isBlacklisted, err := cache.IsTokenBlacklisted(c.Request.Context(), claims.ID)

		if err != nil || isBlacklisted {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code": response.ErrUnauthorized,
				"err":  "Token has been revoked", "description": "",
			})
			return
		}

		ctx := context.WithValue(c.Request.Context(), "subjectUUID", claims.Subject)
		ctx = context.WithValue(ctx, "accessToken", jwtToken)

		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}
