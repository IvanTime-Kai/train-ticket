package model

// @name CreateTrainRequest
type CreateTrainRequest struct {
	Name       string `json:"name"        binding:"required"`
	TotalSeats int    `json:"total_seats" binding:"required,min=1"`
}

// @name CreateSeatRequest
type CreateSeatRequest struct {
	SeatNumber string  `json:"seat_number" binding:"required"`
	Class      string  `json:"class"       binding:"required,oneof=economy business vip"`
	Price      float64 `json:"price"       binding:"required,min=0"`
}

// @name TrainResponse
type TrainResponse struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	TotalSeats int    `json:"total_seats"`
	Status     int    `json:"status"`
}

// @name SeatResponse
type SeatResponse struct {
	ID         string  `json:"id"`
	TrainID    string  `json:"train_id"`
	SeatNumber string  `json:"seat_number"`
	Class      string  `json:"class"`
	Price      float64 `json:"price"`
}
