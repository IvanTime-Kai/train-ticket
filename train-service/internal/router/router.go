package router

type RouterGroup struct {
	StationRouter StationRouter
	TrainRouter   TrainRouter
	TripRouter    TripRouter
}

var RouterGroupApp = new(RouterGroup)
