package router

import (
	"github.com/gin-gonic/gin"
	"github.com/leminhthai/train-ticket/user-service/internal/handler"
	"github.com/leminhthai/train-ticket/user-service/internal/middleware"
)

type UserRouter struct{}

func (ur *UserRouter) InitUserRouter(Router *gin.RouterGroup) {
	userPrivate := Router.Group("/user")
	userPrivate.Use(middleware.AuthMiddleware())
	{
		userPrivate.GET("/profile", handler.User.GetByID)
		userPrivate.PUT("/profile", handler.User.Update)
	}
}
