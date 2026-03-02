package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/leminhthai/train-ticket/user-service/internal/model"
	"github.com/leminhthai/train-ticket/user-service/internal/service"
	"github.com/leminhthai/train-ticket/user-service/pkg/response"
)

type UserHandler struct {
	us service.UserService
}

func NewUserHandler(us service.UserService) *UserHandler {
	return &UserHandler{
		us: us,
	}
}

func (uh *UserHandler) Login(c *gin.Context) {
	var req model.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorResponse(c, response.ErrCodeParamsInvalid, err.Error())
		return
	}

	user, err := uh.us.Login(c.Request.Context(), &req)

	if err != nil {
		response.ErrorResponse(c, response.ErrCodeParamsInvalid, err.Error())
		return
	}

	response.SuccessResponse(c, response.ErrCodeSuccess, user)
}

func (uh *UserHandler) Register(c *gin.Context) {
	var req model.RegisterRequest

	if err := c.ShouldBindJSON(c.Request.Context()); err != nil {
		response.ErrorResponse(c, response.ErrCodeParamsInvalid, err.Error())
		return
	}

	user, err := uh.us.Register(c.Request.Context(), &req)

	if err != nil {
		response.ErrorResponse(c, response.ErrCodeParamsInvalid, err.Error())
		return
	}

	response.SuccessResponse(c, response.ErrCodeSuccess, user)
}

func (uh *UserHandler) GetByID(c *gin.Context) {
	userID := c.Request.Context().Value("subjectUUID").(string)

	user, err := uh.us.GetByID(c.Request.Context(), userID)

	if err != nil {
		response.ErrorResponse(c, response.ErrUserNotFound, err.Error())
		return
	}

	response.SuccessResponse(c, response.ErrCodeSuccess, user)
}

func (uh *UserHandler) Update(c *gin.Context) {
	userID := c.Request.Context().Value("subjectUUID").(string)

	var req model.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorResponse(c, response.ErrCodeParamsInvalid, err.Error())
		return
	}

	if err := uh.us.Update(c.Request.Context(), userID, &req); err != nil {
		response.ErrorResponse(c, response.ErrCodeParamsInvalid, err.Error())
		return
	}

	response.SuccessResponse(c, response.ErrCodeSuccess, nil)
}
