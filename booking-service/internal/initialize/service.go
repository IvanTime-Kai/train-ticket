package initialize

import (
	"github.com/leminhthai/train-ticket/booking-service/global"
	"github.com/leminhthai/train-ticket/booking-service/internal/handler"
	"github.com/leminhthai/train-ticket/booking-service/pkg/wire"
	"go.uber.org/zap"
)

func InitServices() {
	initBookingServices()
}

func initBookingServices() {
	app, err := wire.InitializeApp(global.Mdb, global.Config.GRPC)

	if err != nil {
		global.Logger.Fatal("wire InitializeApp failed", zap.Error(err))
	}

	handler.InitBooking(app.BookingService)
	handler.InitPayment(app.PaymentService)
}
