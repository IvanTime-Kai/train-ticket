package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/leminhthai/train-ticket/booking-service/internal/model"
	"github.com/leminhthai/train-ticket/booking-service/internal/service"
	"github.com/leminhthai/train-ticket/booking-service/pkg/response"
)

var Booking = new(BookingHandler)

type BookingHandler struct {
	bs service.BookingService
}

func InitBooking(bs service.BookingService) {
	Booking.bs = bs
}

// HoldSeat godoc
// @Summary     Giữ ghế tạm thời
// @Description Giữ ghế trong 5 phút trước khi tạo booking
// @Tags        bookings
// @Accept      json
// @Produce     json
// @Security    BearerAuth
// @Param       request body model.HoldSeatRequest true "Hold Seat Request"
// @Success     200 {object} response.ResponseData
// @Failure     400 {object} response.ResponseData
// @Failure     401 {object} response.ResponseData
// @Router      /bookings/hold [post]
func (bh *BookingHandler) HoldSeat(c *gin.Context) {
	var req model.HoldSeatRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorResponse(c, response.ErrCodeParamsInvalid, err.Error())
		return
	}

	userID := c.GetString("subjectUUID")
	result, err := bh.bs.HoldSeat(c.Request.Context(), userID, &req)

	if err != nil {
		response.ErrorResponse(c, response.ErrCodeParamsInvalid, err.Error())
		return
	}

	response.SuccessResponse(c, response.ErrCodeSuccess, result)
}

// CreateBooking godoc
// @Summary     Tạo booking
// @Description Tạo booking sau khi đã giữ ghế
// @Tags        bookings
// @Accept      json
// @Produce     json
// @Security    BearerAuth
// @Param       request body model.CreateBookingRequest true "Create Booking Request"
// @Success     200 {object} response.ResponseData
// @Failure     400 {object} response.ResponseData
// @Failure     401 {object} response.ResponseData
// @Router      /bookings [post]
func (bh *BookingHandler) CreateBooking(c *gin.Context) {
	var req model.CreateBookingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorResponse(c, response.ErrCodeParamsInvalid, err.Error())
		return
	}

	userID := c.GetString("subjectUUID")

	result, err := bh.bs.CreateBooking(c.Request.Context(), userID, &req)
	if err != nil {
		response.ErrorResponse(c, response.ErrCodeParamsInvalid, err.Error())
		return
	}

	response.SuccessResponse(c, response.ErrCodeSuccess, result)
}

// GetBookingByID godoc
// @Summary     Chi tiết booking
// @Description Lấy thông tin booking theo ID
// @Tags        bookings
// @Produce     json
// @Security    BearerAuth
// @Param       id path string true "Booking ID"
// @Success     200 {object} response.ResponseData
// @Failure     404 {object} response.ResponseData
// @Router      /bookings/{id} [get]
func (bh *BookingHandler) GetBookingByID(c *gin.Context) {
	bookingID := c.Param("id")
	userID := c.GetString("subjectUUID")

	result, err := bh.bs.GetBookingByID(c.Request.Context(), userID, bookingID)
	if err != nil {
		response.ErrorResponse(c, response.ErrNotFound, err.Error())
		return
	}

	response.SuccessResponse(c, response.ErrCodeSuccess, result)
}

// GetMyBookings godoc
// @Summary     Danh sách booking của tôi
// @Description Lấy tất cả booking của user đang đăng nhập
// @Tags        bookings
// @Produce     json
// @Security    BearerAuth
// @Success     200 {object} response.ResponseData
// @Router      /bookings/me [get]
func (bh *BookingHandler) GetMyBookings(c *gin.Context) {
	userID := c.GetString("subjectUUID")

	result, err := bh.bs.GetMyBookings(c.Request.Context(), userID)
	if err != nil {
		response.ErrorResponse(c, response.ErrInternalServer, err.Error())
		return
	}

	response.SuccessResponse(c, response.ErrCodeSuccess, result)
}

// CancelBooking godoc
// @Summary     Huỷ booking
// @Description Huỷ booking theo ID
// @Tags        bookings
// @Produce     json
// @Security    BearerAuth
// @Param       id path string true "Booking ID"
// @Success     200 {object} response.ResponseData
// @Failure     400 {object} response.ResponseData
// @Router      /bookings/{id} [delete]
func (bh *BookingHandler) CancelBooking(c *gin.Context) {
	bookingID := c.Param("id")
	userID := c.GetString("subjectUUID")

	if err := bh.bs.CancelBooking(c.Request.Context(), userID, bookingID); err != nil {
		response.ErrorResponse(c, response.ErrCodeParamsInvalid, err.Error())
		return
	}

	response.SuccessResponse(c, response.ErrCodeSuccess, nil)
}
