package router

import (
	"github.com/gin-gonic/gin"
	"github.com/leminhthai/train-ticket/booking-service/internal/handler"
)

type BookingRouter struct{}

func (br *BookingRouter) InitBookingRouter(Router *gin.RouterGroup) {

	// PRIVATE
	bookingPrivate := Router.Group("/bookings")
	{
		bookingPrivate.POST("/hold", handler.Booking.HoldSeat)
		bookingPrivate.POST("", handler.Booking.CreateBooking)
		bookingPrivate.GET("/me", handler.Booking.GetMyBookings)
		bookingPrivate.GET("/:id", handler.Booking.GetBookingByID)
		bookingPrivate.DELETE("/:id", handler.Booking.CancelBooking)
	}
}
