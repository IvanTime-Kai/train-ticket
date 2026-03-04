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

func (ah *AuthHandler) Logout(c *gin.Context) {
	userId, _ := c.Request.Context().Value("subjectUUID").(string)

	if userId == "" {
		response.ErrorResponse(c, response.ErrUnauthorized, "")
		return
	}

	accessToken, _ := c.Request.Context().Value("accessToken").(string)

	if err := ah.us.Logout(c.Request.Context(), userId, accessToken); err != nil {
		response.ErrorResponse(c, response.ErrInternalServer, err.Error())
		return
	}

	response.SuccessResponse(c, response.ErrCodeSuccess, nil)
}

func (ah *AuthHandler) RefreshToken(c *gin.Context) {
	var req model.RefreshTokenRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorResponse(c, response.ErrCodeParamsInvalid, err.Error())
		return
	}

	accessToken, err := ah.us.RefreshToken(c.Request.Context(), req.RefreshToken)

	if err != nil {
		response.ErrorResponse(c, response.ErrUnauthorized, err.Error())
	}

	response.SuccessResponse(c, response.ErrCodeSuccess, gin.H{
		"access_token": accessToken,
	})

}

func (ah *AuthHandler) ForgotPassword(c *gin.Context) {
	var req model.ForgotPasswordRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorResponse(c, response.ErrCodeParamsInvalid, err.Error())
		return
	}

	if err := ah.us.ForgotPassword(c.Request.Context(), &req); err != nil {
		response.ErrorResponse(c, response.ErrInternalServer, err.Error())
		return
	}

	response.SuccessResponse(c, response.ErrCodeSuccess, nil)
}

func (ah *AuthHandler) VerifyOTP(c *gin.Context) {
	var req model.VerifyOTPRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorResponse(c, response.ErrCodeParamsInvalid, err.Error())
		return
	}

	reset_token, err := ah.us.VerifyOTP(c.Request.Context(), &req)
	if err != nil {
		response.ErrorResponse(c, response.ErrCodeParamsInvalid, err.Error())
	}

	response.SuccessResponse(c, response.ErrCodeSuccess, gin.H{
		"reset_token": reset_token,
	})
}

func (uh *AuthHandler) ResetPassword(c *gin.Context) {
	var req model.ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorResponse(c, response.ErrCodeParamsInvalid, err.Error())
		return
	}

	if err := uh.us.ResetPassword(c.Request.Context(), &req); err != nil {
		response.ErrorResponse(c, response.ErrCodeParamsInvalid, err.Error())
		return
	}

	response.SuccessResponse(c, response.ErrCodeSuccess, nil)
}

func (uh *UserHandler) ChangePassword(c *gin.Context) {
	userID, _ := c.Request.Context().Value("subjectUUID").(string)
	if userID == "" {
		response.ErrorResponse(c, response.ErrUnauthorized, "")
		return
	}

	var req model.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorResponse(c, response.ErrCodeParamsInvalid, err.Error())
		return
	}

	if err := uh.us.ChangePassword(c.Request.Context(), userID, &req); err != nil {
		response.ErrorResponse(c, response.ErrCodeParamsInvalid, err.Error())
		return
	}

	response.SuccessResponse(c, response.ErrCodeSuccess, nil)
}