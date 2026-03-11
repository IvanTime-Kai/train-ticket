package router

import (
	"github.com/gin-gonic/gin"
	"github.com/leminhthai/train-ticket/booking-service/internal/handler"
	"github.com/leminhthai/train-ticket/booking-service/internal/middleware"
)

type PaymentRouter struct{}

func (pr *PaymentRouter) InitPaymentRouter(Router *gin.RouterGroup) {

	// PRIVATE
	paymentPrivate := Router.Group("/payments")
	paymentPrivate.Use(middleware.AuthMiddleware())
	{
		paymentPrivate.POST("", handler.Payment.CreatePayment)
	}
}
