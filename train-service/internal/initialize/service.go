package initialize

import (
	"github.com/leminhthai/train-ticket/train-service/global"
	"github.com/leminhthai/train-ticket/train-service/internal/handler"
	"github.com/leminhthai/train-ticket/train-service/pkg/wire"
)

func InitServices() {
	initUserServices()
}

func initUserServices() {
	app := wire.InitializeApp(global.Mdb)

	handler.InitStation(app.StationService)
	handler.InitTrain(app.TrainService)
	handler.InitTrip(app.TripService)
}
