package initialize

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/leminhthai/train-ticket/train-service/global"
	"github.com/leminhthai/train-ticket/train-service/internal/router"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func InitRouter() *gin.Engine {
	var r *gin.Engine

	if global.Config.Server.Mode == "dev" {
		gin.SetMode(gin.DebugMode)
		gin.ForceConsoleColor()
		r = gin.Default()
	} else {
		gin.SetMode(gin.ReleaseMode)
		r = gin.New()
	}

	r.GET("/health-check", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "User Service is running"})
	})

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	MainGroup := r.Group("/api/v1")
	stationRouter := router.RouterGroupApp.StationRouter
	trainRouter := router.RouterGroupApp.TrainRouter
	tripRouter := router.RouterGroupApp.TripRouter
	{
		stationRouter.InitStationRouter(MainGroup)
		trainRouter.InitTrainRouterRouter(MainGroup)
		tripRouter.InitTripRouterRouter(MainGroup)
	}

	return r
}
