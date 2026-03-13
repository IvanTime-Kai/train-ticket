package initialize

import (
	"github.com/leminhthai/train-ticket/booking-service/global"
	"github.com/leminhthai/train-ticket/booking-service/internal/worker"
	"github.com/leminhthai/train-ticket/booking-service/pkg/wire"
)

func InitWorker(app *wire.App) {
	w := worker.NewExpiryWorker(
		app.BookingRepository,
		global.Logger.Logger,
	)

	w.Start()
}
