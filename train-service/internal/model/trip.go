package model

const (
	TripStatusScheduled = 1 // chờ khởi hành
	TripStatusDeparted  = 2 // đã khởi hành
	TripStatusArrived   = 3 // đã đến nơi
	TripStatusCancelled = 0 // đã huỷ
)

// Train Status Constants
const (
	TrainStatusActive   = 1 // đang hoạt động
	TrainStatusInactive = 0 // ngừng hoạt động
)

// Seat Class Constants
const (
	SeatClassEconomy  = "economy"
	SeatClassBusiness = "business"
	SeatClassVIP      = "vip"
)

// @name CreateTripRequest
type CreateTripRequest struct {
	TrainID       string `json:"train_id"       binding:"required"`
	RouteID       string `json:"route_id"       binding:"required"`
	DepartureTime string `json:"departure_time" binding:"required"` // format: 2006-01-02 15:04:05
	ArrivalTime   string `json:"arrival_time"   binding:"required"`
}

// @name SearchTripRequest
type SearchTripRequest struct {
	From string `form:"from" binding:"required"` // station code VD: HAN
	To   string `form:"to"   binding:"required"` // station code VD: SGN
	Date string `form:"date" binding:"required"` // format: 2006-01-02
}

// @name TripResponse
type TripResponse struct {
	ID            string `json:"id"`
	TrainID       string `json:"train_id"`
	RouteID       string `json:"route_id"`
	DepartureTime string `json:"departure_time"`
	ArrivalTime   string `json:"arrival_time"`
	Status        int    `json:"status"`
}

// @name CreateRouteRequest
type CreateRouteRequest struct {
	OriginStationID      string `json:"origin_station_id"      binding:"required"`
	DestinationStationID string `json:"destination_station_id" binding:"required"`
	DistanceKm           int    `json:"distance_km"`
}
