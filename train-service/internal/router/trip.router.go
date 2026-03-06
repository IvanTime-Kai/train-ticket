package router

import (
	"github.com/gin-gonic/gin"
	"github.com/leminhthai/train-ticket/train-service/internal/handler"
)

type TripRouter struct{}

func (tr *TripRouter) InitTripRouterRouter(Router *gin.RouterGroup) {

	// PUBLIC
	tripPublic := Router.Group("/trips")
	{
		tripPublic.GET("/search", handler.Trip.Search)
		tripPublic.GET("/:id", handler.Trip.GetByID)
	}

	// PRIVATE
	tripPrivate := Router.Group("/trips")
	{
		tripPrivate.POST("", handler.Trip.Create)
		tripPrivate.POST("/:id/cancel", handler.Trip.Cancel)
	}

	// Routes
	routeGroup := Router.Group("/routes")
	{
		routeGroup.POST("", handler.Trip.CreateRoute)
	}
}
