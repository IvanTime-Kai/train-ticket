package initialize

import (
	"github.com/leminhthai/train-ticket/booking-service/global"
	"github.com/leminhthai/train-ticket/booking-service/internal/handler"
	"github.com/leminhthai/train-ticket/booking-service/pkg/wire"
)

func InitServices() {
	initBookingServices()
}

func initBookingServices() {
	app := wire.InitializeApp(global.Mdb)

	handler.InitBooking(app.BookingService)
	handler.InitPayment(app.PaymentService)
}
