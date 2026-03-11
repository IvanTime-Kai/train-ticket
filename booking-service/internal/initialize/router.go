package initialize

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/leminhthai/train-ticket/booking-service/global"
	"github.com/leminhthai/train-ticket/booking-service/internal/router"
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
		c.JSON(http.StatusOK, gin.H{"message": "Booking Service is running"})
	})

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	MainGroup := r.Group("/api/v1")
	bookingRouter := router.RouterGroupApp.BookingRouter
	paymentRouter := router.RouterGroupApp.PaymentRouter
	{
		bookingRouter.InitBookingRouter(MainGroup)
		paymentRouter.InitPaymentRouter(MainGroup)
	}

	return r
}
