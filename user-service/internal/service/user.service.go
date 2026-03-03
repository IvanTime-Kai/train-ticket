package service

import (
	"context"
	"fmt"

	db "github.com/leminhthai/train-ticket/user-service/db/generated"
	"github.com/leminhthai/train-ticket/user-service/internal/model"
	"github.com/leminhthai/train-ticket/user-service/internal/repository"
	"github.com/leminhthai/train-ticket/user-service/internal/utils/auth"
	"github.com/leminhthai/train-ticket/user-service/internal/utils/cache"
	"github.com/leminhthai/train-ticket/user-service/internal/utils/crypto"
)

type UserService interface {
	Register(ctx context.Context, req *model.RegisterRequest) (db.User, error)
	Login(ctx context.Context, req *model.LoginRequest) (*model.LoginResponse, error)
	GetByID(ctx context.Context, id string) (db.User, error)
	Update(ctx context.Context, id string, req *model.UpdateUserRequest) error
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
