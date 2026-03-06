package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/leminhthai/train-ticket/train-service/internal/model"
	"github.com/leminhthai/train-ticket/train-service/internal/service"
	"github.com/leminhthai/train-ticket/train-service/pkg/response"
)

var Trip = new(TripHandler)

type TripHandler struct {
	ts service.TripService
}

func InitTrip(ts service.TripService) {
	Trip.ts = ts
}

// @Summary     Tạo chuyến tàu
// @Description Tạo chuyến tàu mới
// @Tags        trips
// @Accept      json
// @Produce     json
// @Param       request body model.CreateTripRequest true "Create Trip Request"
// @Success     200 {object} response.ResponseData
// @Failure     400 {object} response.ResponseData
// @Router      /trips [post]
func (h *TripHandler) Create(c *gin.Context) {
	var req model.CreateTripRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorResponse(c, response.ErrCodeParamsInvalid, err.Error())
		return
	}

	trip, err := h.ts.Create(c.Request.Context(), &req)
	if err != nil {
		response.ErrorResponse(c, response.ErrCodeParamsInvalid, err.Error())
		return
	}

	response.SuccessResponse(c, response.ErrCodeSuccess, trip)
}

// @Summary     Chi tiết chuyến tàu
// @Description Lấy thông tin chuyến tàu theo ID
// @Tags        trips
// @Produce     json
// @Param       id path string true "Trip ID"
// @Success     200 {object} response.ResponseData
// @Failure     404 {object} response.ResponseData
// @Router      /trips/{id} [get]
func (h *TripHandler) GetByID(c *gin.Context) {
	id := c.Param("id")

	trip, err := h.ts.GetByID(c.Request.Context(), id)
	if err != nil {
		response.ErrorResponse(c, response.ErrNotFound, err.Error())
		return
	}

	response.SuccessResponse(c, response.ErrCodeSuccess, trip)
}

// @Summary     Tìm kiếm chuyến tàu
// @Description Tìm chuyến tàu theo ga và ngày
// @Tags        trips
// @Produce     json
// @Param       from  query string true "Mã ga xuất phát (VD: HAN)"
// @Param       to    query string true "Mã ga đến (VD: SGN)"
// @Param       date  query string true "Ngày khởi hành (YYYY-MM-DD)"
// @Success     200 {object} response.ResponseData
// @Failure     400 {object} response.ResponseData
// @Router      /trips/search [get]
func (h *TripHandler) Search(c *gin.Context) {
	var req model.SearchTripRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.ErrorResponse(c, response.ErrCodeParamsInvalid, err.Error())
		return
	}

	trips, err := h.ts.Search(c.Request.Context(), &req)
	if err != nil {
		response.ErrorResponse(c, response.ErrCodeParamsInvalid, err.Error())
		return
	}

	response.SuccessResponse(c, response.ErrCodeSuccess, trips)
}

// @Summary     Huỷ chuyến tàu
// @Description Huỷ chuyến tàu theo ID
// @Tags        trips
// @Produce     json
// @Param       id path string true "Trip ID"
// @Success     200 {object} response.ResponseData
// @Failure     404 {object} response.ResponseData
// @Router      /trips/{id}/cancel [post]
func (h *TripHandler) Cancel(c *gin.Context) {
	id := c.Param("id")

	if err := h.ts.Cancel(c.Request.Context(), id); err != nil {
		response.ErrorResponse(c, response.ErrNotFound, err.Error())
		return
	}

	response.SuccessResponse(c, response.ErrCodeSuccess, nil)
}

// @Summary     Tạo tuyến đường
// @Description Tạo tuyến đường mới giữa 2 ga
// @Tags        routes
// @Accept      json
// @Produce     json
// @Param       request body model.CreateRouteRequest true "Create Route Request"
// @Success     200 {object} response.ResponseData
// @Failure     400 {object} response.ResponseData
// @Router      /routes [post]
func (h *TripHandler) CreateRoute(c *gin.Context) {
	var req model.CreateRouteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorResponse(c, response.ErrCodeParamsInvalid, err.Error())
		return
	}

	route, err := h.ts.CreateRoute(c.Request.Context(), &req)
	if err != nil {
		response.ErrorResponse(c, response.ErrCodeParamsInvalid, err.Error())
		return
	}

	response.SuccessResponse(c, response.ErrCodeSuccess, route)
}
