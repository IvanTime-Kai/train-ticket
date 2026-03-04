package service

import (
	"context"
	"fmt"
	"time"

	db "github.com/leminhthai/train-ticket/user-service/db/generated"
	"github.com/leminhthai/train-ticket/user-service/internal/model"
	"github.com/leminhthai/train-ticket/user-service/internal/repository"
	"github.com/leminhthai/train-ticket/user-service/internal/utils/auth"
	"github.com/leminhthai/train-ticket/user-service/internal/utils/cache"
	"github.com/leminhthai/train-ticket/user-service/internal/utils/crypto"
	"github.com/leminhthai/train-ticket/user-service/pkg/email"
)

type UserService interface {
	Register(ctx context.Context, req *model.RegisterRequest) (db.User, error)
	Login(ctx context.Context, req *model.LoginRequest) (*model.LoginResponse, error)
	GetByID(ctx context.Context, id string) (db.User, error)
	Update(ctx context.Context, id string, req *model.UpdateUserRequest) error
	Logout(ctx context.Context, userId, accessToken string) error
	RefreshToken(ctx context.Context, refreshToken string) (string, error)
	ForgotPassword(ctx context.Context, req *model.ForgotPasswordRequest) error
	VerifyOTP(ctx context.Context, req *model.VerifyOTPRequest) (string, error)
	ResetPassword(ctx context.Context, req *model.ResetPasswordRequest) error
	ChangePassword(ctx context.Context, userId string, req *model.ChangePasswordRequest) error
}

type userService struct {
	repo repository.UserRepository
}

// NewUserService returns a UserService that uses the given UserRepository.
func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo: repo}
}

func (us *userService) Register(ctx context.Context, req *model.RegisterRequest) (db.User, error) {
	existingUser, err := us.repo.GetByEmail(ctx, req.Email)
	if err == nil && existingUser.ID != "" {
		return db.User{}, fmt.Errorf("email already exists")
	}

	hashedPassword, err := crypto.HashPassword(req.Password)
	if err != nil {
		return db.User{}, err
	}

	return us.repo.Create(ctx, req, hashedPassword)
}
func (us *userService) Login(ctx context.Context, req *model.LoginRequest) (*model.LoginResponse, error) {
	// check email exist
	existingUser, err := us.repo.GetByEmail(ctx, req.Email)

	if err != nil {
		return nil, fmt.Errorf("invalid email or password")
	}

	// compare password
	if !crypto.MatchingPassword(existingUser.Password, req.Password) {
		return nil, fmt.Errorf("invalid email or password")
	}

	// generated access token
	accessToken, err := auth.GenerateAccessToken(existingUser.ID)
	if err != nil {
		return nil, err
	}

	// generated refresh token
	refreshToken, err := auth.GenerateRefreshToken(existingUser.ID)
	if err != nil {
		return nil, err
	}

	// save refresh token
	if err := cache.SaveRefreshToken(ctx, existingUser.ID, refreshToken); err != nil {
		return nil, err
	}

	// create session
	if err := us.repo.CreateSession(ctx, existingUser.ID, req.Device, req.IPAddress); err != nil {
		return nil, err
	}

	if err := us.repo.UpdateLastLogin(ctx, existingUser.ID); err != nil {
		return nil, err
	}

	return &model.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User: model.UserResponse{
			ID:         existingUser.ID,
			Email:      existingUser.Email,
			FullName:   existingUser.FullName,
			Phone:      existingUser.Phone.String,
			Role:       existingUser.Role,
			IsVerified: existingUser.IsVerified,
		},
	}, nil

}
func (us *userService) GetByID(ctx context.Context, id string) (db.User, error) {
	user, err := us.repo.GetByID(ctx, id)
	if err != nil {
		return db.User{}, fmt.Errorf("user not found")
	}
	return user, nil
}
func (us *userService) Update(ctx context.Context, id string, req *model.UpdateUserRequest) error {
	return us.repo.Update(ctx, id, req)
}
func (us *userService) Logout(ctx context.Context, userId, accessToken string) error {
	// Parse access token lấy claims
	claims, err := auth.ParseToken(accessToken)
	if err != nil {
		return fmt.Errorf("invalid token")
	}

	// Calculate remaining
	remaining := time.Until(claims.ExpiresAt.Time)

	if remaining > 0 {
		if err := cache.BlacklistToken(ctx, claims.ID, remaining); err != nil {
			return err
		}
	}

	// remove refresh token
	if err := cache.DeleteRefreshToken(ctx, userId); err != nil {
		return err
	}

	return us.repo.LogoutAllSessions(ctx, userId)
}
func (us *userService) RefreshToken(ctx context.Context, refreshToken string) (string, error) {
	claims, err := auth.ParseToken(refreshToken)

	if err != nil {
		return "", fmt.Errorf("invalid refresh token")
	}

	if claims.TokenType != auth.TokenTypeRefresh {
		return "", fmt.Errorf("invalid token type")
	}

	userId := claims.Subject

	savedToken, err := cache.GetRefreshToken(ctx, userId)
	if err != nil || savedToken != refreshToken {
		return "", fmt.Errorf("session expired, please login again")
	}

	accessToken, err := auth.GenerateAccessToken(userId)

	if err != nil {
		return "", err
	}

	return accessToken, nil
}
func (us *userService) ForgotPassword(ctx context.Context, req *model.ForgotPasswordRequest) error {
	// check email exist
	_, err := us.repo.GetByEmail(ctx, req.Email)

	if err != nil {
		return nil
	}

	// generate otp
	otp, err := crypto.GenerateOTP(6)

	if err != nil {
		return err
	}

	// set otp into redis
	if err := cache.SaveOTP(ctx, req.Email, otp); err != nil {
		return err
	}

	// send email
	return email.SendOTPEmail(req.Email, otp)
}
func (us *userService) VerifyOTP(ctx context.Context, req *model.VerifyOTPRequest) (string, error) {
	// get otp
	storedOtp, err := cache.GetOTP(ctx, req.Email)
	if err != nil {
		return "", fmt.Errorf("OTP expired or invalid")
	}

	// compare OTP
	if storedOtp != req.OTP {
		return "", fmt.Errorf("OTP incorrect")
	}

	// delete OTP in redis
	if err := cache.DeleteOTP(ctx, req.Email); err != nil {
		return "", err
	}

	// generate reset token
	resetToken, err := auth.GenerateAccessToken(req.Email)
	if err != nil {
		return "", err
	}

	// set reset token into redis
	if err := cache.SaveResetToken(ctx, req.Email, resetToken); err != nil {
		return "", err
	}

	return resetToken, nil
}
func (us *userService) ResetPassword(ctx context.Context, req *model.ResetPasswordRequest) error {
	// get email
	claims, err := auth.ParseToken(req.ResetToken)
	if err != nil {
		return err
	}
	email := claims.Subject

	// check reset token
	storedResetToken, err := cache.GetResetToken(ctx, email)
	if err != nil || storedResetToken != req.ResetToken {
		return fmt.Errorf("reset token expired or invalid")
	}

	// get userId
	user, err := us.repo.GetByEmail(ctx, email)
	if err != nil {
		return fmt.Errorf("user not found")
	}

	// hashed password
	hashedPassword, err := crypto.HashPassword(req.NewPassword)
	if err != nil {
		return err
	}

	// update new password into DB
	if err := us.repo.UpdatePassword(ctx, user.ID, hashedPassword); err != nil {
		return err
	}

	// delete reset token in redis
	if err := cache.DeleteResetToken(ctx, email); err != nil {
		return err
	}

	// delete refresh token in redis
	return cache.DeleteRefreshToken(ctx, user.ID)
}
func (us *userService) ChangePassword(ctx context.Context, userId string, req *model.ChangePasswordRequest) error {
	// get user
	existingUser, err := us.GetByID(ctx, userId)

	if err != nil {
		return fmt.Errorf("user not found")
	}

	// compare password
	if !crypto.MatchingPassword(existingUser.Password, req.OldPassword) {
		return fmt.Errorf("old password incorrect")
	}

	// hashed new password
	hashedPassword, err := crypto.HashPassword(req.NewPassword)
	if err != nil {
		return err
	}

	// update password in db
	if err := us.repo.UpdatePassword(ctx, existingUser.ID, hashedPassword); err != nil {
		return err
	}

	// delete refresh token
	if err := cache.DeleteRefreshToken(ctx, existingUser.ID); err != nil {
		return err
	}

	// added current access token into black list
	accessToken := ctx.Value("accessToken").(string)
	if accessToken != "" {
		claims, err := auth.ParseToken(accessToken)
		if err == nil {
			remaining := time.Until(claims.ExpiresAt.Time)
			if remaining > 0 {
				cache.BlacklistToken(ctx, claims.ID, remaining)
			}
		}
	}

	return nil
}
