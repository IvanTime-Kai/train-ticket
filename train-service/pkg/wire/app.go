package wire

import "github.com/leminhthai/train-ticket/train-service/internal/service"

type App struct {
	StationService service.StationService
	TrainService   service.TrainService
	TripService    service.TripService
}
