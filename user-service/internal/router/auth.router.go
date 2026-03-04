package router

import (
	"github.com/gin-gonic/gin"
	"github.com/leminhthai/train-ticket/user-service/internal/handler"
	"github.com/leminhthai/train-ticket/user-service/internal/middleware"
)

type AuthRouter struct{}

func (ar *AuthRouter) InitAuthRouter(Router *gin.RouterGroup) {

	// PUBLIC
	authPublic := Router.Group("/auth")
	{
		authPublic.POST("/register", handler.Auth.Register)
		authPublic.POST("/login", handler.Auth.Login)
		authPublic.POST("/refresh-token", handler.Auth.RefreshToken)
		authPublic.POST("/forgot-password", handler.Auth.ForgotPassword)
		authPublic.POST("/verify-otp",     handler.Auth.VerifyOTP)
		authPublic.POST("/reset-password", handler.Auth.ResetPassword)
	}

	// PRIVATE
	authPrivate := Router.Group("/auth")
	authPrivate.Use(middleware.AuthMiddleware())
	{
		authPrivate.POST("/logout", handler.Auth.Logout)
		authPrivate.POST("/change-password", handler.User.ChangePassword)
	}
}
