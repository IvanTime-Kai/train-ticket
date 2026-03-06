package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/leminhthai/train-ticket/train-service/internal/model"
	"github.com/leminhthai/train-ticket/train-service/internal/service"
	"github.com/leminhthai/train-ticket/train-service/pkg/response"
)

var Station = new(StationHandler)

type StationHandler struct {
	ss service.StationService
}

func InitStation(ss service.StationService) {
	Station.ss = ss
}

// @Summary     Tạo ga tàu
// @Description Tạo ga tàu mới
// @Tags        stations
// @Accept      json
// @Produce     json
// @Param       request body model.CreateStationRequest true "Create Station Request"
// @Success     200 {object} response.ResponseData
// @Failure     400 {object} response.ResponseData
// @Router      /stations [post]
func (sh *StationHandler) Create(c *gin.Context) {
	var req model.CreateStationRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorResponse(c, response.ErrCodeParamsInvalid, err.Error())
		return
	}

	station, err := sh.ss.Create(c.Request.Context(), &req)
	if err != nil {
		response.ErrorResponse(c, response.ErrCodeParamsInvalid, err.Error())
		return
	}

	response.SuccessResponse(c, response.ErrCodeSuccess, station)
}

// @Summary     Chi tiết ga tàu
// @Description Lấy thông tin ga tàu theo ID
// @Tags        stations
// @Produce     json
// @Param       id path string true "Station ID"
// @Success     200 {object} response.ResponseData
// @Failure     404 {object} response.ResponseData
// @Router      /stations/{id} [get]
func (sh *StationHandler) GetByID(c *gin.Context) {
	id := c.Param("id")

	station, err := sh.ss.GetById(c.Request.Context(), id)
	if err != nil {
		response.ErrorResponse(c, response.ErrNotFound, err.Error())
		return
	}

	response.SuccessResponse(c, response.ErrCodeSuccess, station)
}

// @Summary     Danh sách ga tàu
// @Description Lấy danh sách tất cả ga tàu
// @Tags        stations
// @Produce     json
// @Success     200 {object} response.ResponseData
// @Router      /stations [get]
func (h *StationHandler) List(c *gin.Context) {
	stations, err := h.ss.List(c.Request.Context())
	if err != nil {
		response.ErrorResponse(c, response.ErrInternalServer, err.Error())
		return
	}

	response.SuccessResponse(c, response.ErrCodeSuccess, stations)
}
