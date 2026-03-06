package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/leminhthai/train-ticket/train-service/internal/model"
	"github.com/leminhthai/train-ticket/train-service/internal/service"
	"github.com/leminhthai/train-ticket/train-service/pkg/response"
)

var Train = new(TrainHandler)

type TrainHandler struct {
	ts service.TrainService
}

func InitTrain(ts service.TrainService) {
	Train.ts = ts
}

// @Summary     Tạo đoàn tàu
// @Description Tạo đoàn tàu mới
// @Tags        trains
// @Accept      json
// @Produce     json
// @Param       request body model.CreateTrainRequest true "Create Train Request"
// @Success     200 {object} response.ResponseData
// @Failure     400 {object} response.ResponseData
// @Router      /trains [post]
func (h *TrainHandler) Create(c *gin.Context) {
	var req model.CreateTrainRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorResponse(c, response.ErrCodeParamsInvalid, err.Error())
		return
	}

	train, err := h.ts.Create(c.Request.Context(), &req)
	if err != nil {
		response.ErrorResponse(c, response.ErrCodeParamsInvalid, err.Error())
		return
	}

	response.SuccessResponse(c, response.ErrCodeSuccess, train)
}

// @Summary     Chi tiết đoàn tàu
// @Description Lấy thông tin đoàn tàu theo ID
// @Tags        trains
// @Produce     json
// @Param       id path string true "Train ID"
// @Success     200 {object} response.ResponseData
// @Failure     404 {object} response.ResponseData
// @Router      /trains/{id} [get]
func (h *TrainHandler) GetByID(c *gin.Context) {
	id := c.Param("id")

	train, err := h.ts.GetByID(c.Request.Context(), id)
	if err != nil {
		response.ErrorResponse(c, response.ErrNotFound, err.Error())
		return
	}

	response.SuccessResponse(c, response.ErrCodeSuccess, train)
}

// @Summary     Danh sách đoàn tàu
// @Description Lấy danh sách tất cả đoàn tàu
// @Tags        trains
// @Produce     json
// @Success     200 {object} response.ResponseData
// @Router      /trains [get]
func (h *TrainHandler) List(c *gin.Context) {
	trains, err := h.ts.List(c.Request.Context())
	if err != nil {
		response.ErrorResponse(c, response.ErrInternalServer, err.Error())
		return
	}

	response.SuccessResponse(c, response.ErrCodeSuccess, trains)
}

// @Summary     Vô hiệu hoá đoàn tàu
// @Description Đổi trạng thái đoàn tàu thành inactive
// @Tags        trains
// @Produce     json
// @Param       id path string true "Train ID"
// @Success     200 {object} response.ResponseData
// @Failure     404 {object} response.ResponseData
// @Router      /trains/{id} [delete]
func (h *TrainHandler) DeActive(c *gin.Context) {
	id := c.Param("id")

	if err := h.ts.DeActive(c.Request.Context(), id); err != nil {
		response.ErrorResponse(c, response.ErrNotFound, err.Error())
		return
	}

	response.SuccessResponse(c, response.ErrCodeSuccess, nil)
}

// @Summary     Thêm ghế
// @Description Thêm ghế mới vào đoàn tàu
// @Tags        trains
// @Accept      json
// @Produce     json
// @Param       id path string true "Train ID"
// @Param       request body model.CreateSeatRequest true "Create Seat Request"
// @Success     200 {object} response.ResponseData
// @Failure     400 {object} response.ResponseData
// @Router      /trains/{id}/seats [post]
func (h *TrainHandler) AddSeat(c *gin.Context) {
	trainID := c.Param("id")

	var req model.CreateSeatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorResponse(c, response.ErrCodeParamsInvalid, err.Error())
		return
	}

	seat, err := h.ts.AddSeat(c.Request.Context(), trainID, &req)
	if err != nil {
		response.ErrorResponse(c, response.ErrCodeParamsInvalid, err.Error())
		return
	}

	response.SuccessResponse(c, response.ErrCodeSuccess, seat)
}

// @Summary     Danh sách ghế
// @Description Lấy danh sách ghế của đoàn tàu
// @Tags        trains
// @Produce     json
// @Param       id path string true "Train ID"
// @Success     200 {object} response.ResponseData
// @Router      /trains/{id}/seats [get]
func (h *TrainHandler) ListSeat(c *gin.Context) {
	trainID := c.Param("id")

	seats, err := h.ts.ListSeat(c.Request.Context(), trainID)
	if err != nil {
		response.ErrorResponse(c, response.ErrInternalServer, err.Error())
		return
	}

	response.SuccessResponse(c, response.ErrCodeSuccess, seats)
}
