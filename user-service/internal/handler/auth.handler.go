package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/leminhthai/train-ticket/user-service/internal/model"
	"github.com/leminhthai/train-ticket/user-service/internal/service"
	"github.com/leminhthai/train-ticket/user-service/pkg/response"
)

var Auth = new(AuthHandler)

type AuthHandler struct {
	us service.UserService
}

// InitAuth inject UserService cho auth (login, register). Gọi trong InitServices.
func InitAuth(us service.UserService) {
	Auth.us = us
}

func (ah *AuthHandler) Login(c *gin.Context) {
	var req model.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorResponse(c, response.ErrCodeParamsInvalid, err.Error())
		return
	}

	req.IPAddress = c.ClientIP()
	req.Device = c.GetHeader("User-Agent")

	user, err := ah.us.Login(c.Request.Context(), &req)
	if err != nil {
		response.ErrorResponse(c, response.ErrCodeParamsInvalid, err.Error())
		return
	}

	response.SuccessResponse(c, response.ErrCodeSuccess, user)
}

func (ah *AuthHandler) Register(c *gin.Context) {
	var req model.RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorResponse(c, response.ErrCodeParamsInvalid, err.Error())
		return
	}

	user, err := ah.us.Register(c.Request.Context(), &req)
	if err != nil {
		response.ErrorResponse(c, response.ErrCodeParamsInvalid, err.Error())
		return
	}

	response.SuccessResponse(c, response.ErrCodeSuccess, user)
}
