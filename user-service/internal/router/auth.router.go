package router

import (
	"github.com/gin-gonic/gin"
	"github.com/leminhthai/train-ticket/user-service/internal/handler"
)

type AuthRouter struct{}

func (ar *AuthRouter) InitAuthRouter(Router *gin.RouterGroup) {
	// PUBLIC
	authPublic := Router.Group("/auth")
	{
		authPublic.POST("/register", handler.Auth.Register)
		authPublic.POST("/login", handler.Auth.Login)
	}

	// PRIVATE
}
