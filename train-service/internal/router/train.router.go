package router

import (
	"github.com/gin-gonic/gin"
	"github.com/leminhthai/train-ticket/train-service/internal/handler"
)

type TrainRouter struct{}

func (tr *TrainRouter) InitTrainRouterRouter(Router *gin.RouterGroup) {

	// PUBLIC
	trainPublic := Router.Group("/trains")
	{
		trainPublic.GET("", handler.Train.List)
		trainPublic.GET("/:id", handler.Train.GetByID)
	}

	// PRIVATE
	trainPrivate := Router.Group("/trains")
	{
		trainPrivate.POST("", handler.Train.Create)
		trainPrivate.DELETE("/:id", handler.Train.DeActive)
		trainPrivate.POST("/:id/seats", handler.Train.AddSeat)
		trainPrivate.GET("/:id/seats", handler.Train.ListSeat)
	}
}
