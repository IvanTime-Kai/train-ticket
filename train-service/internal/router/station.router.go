package router

import (
	"github.com/gin-gonic/gin"
	"github.com/leminhthai/train-ticket/train-service/internal/handler"
)

type StationRouter struct{}

func (sr *StationRouter) InitStationRouter(Router *gin.RouterGroup) {

	// PUBLIC
	stationPublic := Router.Group("/stations")
	{
		stationPublic.GET("", handler.Station.List)
		stationPublic.GET("/:id", handler.Station.GetByID)
	}

	// PRIVATE
	stationPrivate := Router.Group("/stations")
	{
		stationPrivate.POST("", handler.Station.Create)
	}
}
