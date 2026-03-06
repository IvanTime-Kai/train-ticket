package model

// @name CreateStationRequest
type CreateStationRequest struct {
	Name string `json:"name" binding:"required"`
	Code string `json:"code" binding:"required"`
	City string `json:"city" binding:"required"`
}

// @name StationResponse
type StationResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Code string `json:"code"`
	City string `json:"city"`
}
