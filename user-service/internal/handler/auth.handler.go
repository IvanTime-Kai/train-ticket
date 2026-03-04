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

// @Summary     Login
// @Description Đăng nhập
// @Tags        auth
// @Accept      json
// @Produce     json
// @Param       request body model.LoginRequest true "Login Request"
// @Success     200 {object} response.ResponseData
// @Failure     400 {object} response.ResponseData
// @Failure     401 {object} response.ResponseData
// @Router      /auth/login [post]
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

// @Summary     Register
// @Description Đăng ký tài khoản mới
// @Tags        auth
// @Accept      json
// @Produce     json
// @Param       request body model.RegisterRequest true "Register Request"
// @Success     201 {object} response.ResponseData
// @Failure     400 {object} response.ResponseData
// @Failure     409 {object} response.ResponseData
// @Router      /auth/register [post]
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

// @Summary     Logout
// @Description Đăng xuất
// @Tags        auth
// @Security    BearerAuth
// @Success     200 {object} response.ResponseData
// @Failure     401 {object} response.ResponseData
// @Router      /auth/logout [post]
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

// @Summary     Refresh Token
// @Description Lấy Access Token mới
// @Tags        auth
// @Accept      json
// @Produce     json
// @Param       request body model.RefreshTokenRequest true "Refresh Token Request"
// @Success     200 {object} response.ResponseData
// @Failure     401 {object} response.ResponseData
// @Router      /auth/refresh-token [post]
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

// @Summary     Forgot Password
// @Description Gửi OTP qua email
// @Tags        auth
// @Accept      json
// @Produce     json
// @Param       request body model.ForgotPasswordRequest true "Forgot Password Request"
// @Success     200 {object} response.ResponseData
// @Failure     400 {object} response.ResponseData
// @Router      /auth/forgot-password [post]
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

// @Summary     Verify OTP
// @Description Xác thực OTP
// @Tags        auth
// @Accept      json
// @Produce     json
// @Param       request body model.VerifyOTPRequest true "Verify OTP Request"
// @Success     200 {object} response.ResponseData
// @Failure     400 {object} response.ResponseData
// @Router      /auth/verify-otp [post]
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

// @Summary     Reset Password
// @Description Đặt lại mật khẩu
// @Tags        auth
// @Accept      json
// @Produce     json
// @Param       request body model.ResetPasswordRequest true "Reset Password Request"
// @Success     200 {object} response.ResponseData
// @Failure     400 {object} response.ResponseData
// @Router      /auth/reset-password [post]
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

// @Summary     Change Password
// @Description Đổi mật khẩu
// @Tags        auth
// @Security    BearerAuth
// @Accept      json
// @Produce     json
// @Param       request body model.ChangePasswordRequest true "Change Password Request"
// @Success     200 {object} response.ResponseData
// @Failure     400 {object} response.ResponseData
// @Failure     401 {object} response.ResponseData
// @Router      /auth/change-password [post]
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