package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/leminhthai/train-ticket/user-service/internal/model"
	"github.com/leminhthai/train-ticket/user-service/internal/service"
	"github.com/leminhthai/train-ticket/user-service/pkg/response"
)

var User = new(UserHandler)

type UserHandler struct {
	us service.UserService
}

// InitUser inject UserService cho user (profile). Gọi trong InitServices.
func InitUser(us service.UserService) {
	User.us = us
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
