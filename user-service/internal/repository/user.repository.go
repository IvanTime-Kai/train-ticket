package repository

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	db "github.com/leminhthai/train-ticket/user-service/db/generated"
	"github.com/leminhthai/train-ticket/user-service/internal/model"
)

type UserRepository interface {
	Create(ctx context.Context, req *model.RegisterRequest, hashedPassword string) (db.User, error)
	GetByEmail(ctx context.Context, email string) (db.User, error)
	GetByID(ctx context.Context, id string) (db.User, error)
	Update(ctx context.Context, id string, req *model.UpdateUserRequest) error
	UpdateLastLogin(ctx context.Context, id string) error
	CreateSession(ctx context.Context, userId, device, ipAddress string) error
	UpdateSessionLogout(ctx context.Context, sessionId string) error
	LogoutAllSessions(ctx context.Context, userId string) error
	UpdatePassword(ctx context.Context, id, hashedPassword string) error
}

type userRepository struct {
	queries *db.Queries
}

func NewUserRepository(queries *db.Queries) UserRepository {
	return &userRepository{queries: queries}
}

func (r *userRepository) Create(ctx context.Context, req *model.RegisterRequest, hashedPassword string) (db.User, error) {
	id := uuid.New().String()

	err := r.queries.CreateUser(ctx, db.CreateUserParams{
		ID:       id,
		Email:    req.Email,
		Password: hashedPassword,
		FullName: req.FullName,
		Phone:    sql.NullString{String: req.Phone, Valid: req.Phone != ""},
		Role:     "customer",
	})
	if err != nil {
		return db.User{}, err
	}

	return r.queries.GetUserByID(ctx, id)
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (db.User, error) {
	return r.queries.GetUserByEmail(ctx, email)
}

func (r *userRepository) GetByID(ctx context.Context, id string) (db.User, error) {
	return r.queries.GetUserByID(ctx, id)
}

func (r *userRepository) Update(ctx context.Context, id string, req *model.UpdateUserRequest) error {
	return r.queries.UpdateUser(ctx, db.UpdateUserParams{
		ID:       id,
		FullName: req.FullName,
		Phone:    sql.NullString{String: req.Phone, Valid: req.Phone != ""},
	})
}

func (r *userRepository) UpdateLastLogin(ctx context.Context, id string) error {
	return r.queries.UpdateLastLogin(ctx, id)
}

func (r *userRepository) CreateSession(ctx context.Context, userId, device, ipAddress string) error {
	return r.queries.CreateSession(ctx, db.CreateSessionParams{
		ID:        uuid.New().String(),
		UserID:    userId,
		Device:    sql.NullString{String: device, Valid: device != ""},
		IpAddress: sql.NullString{String: ipAddress, Valid: ipAddress != ""},
	})
}

func (r *userRepository) UpdateSessionLogout(ctx context.Context, sessionId string) error {
	return r.queries.UpdateSessionLogout(ctx, sessionId)
}

func (r *userRepository) LogoutAllSessions(ctx context.Context, userId string) error {
	return r.queries.LogoutAllSessions(ctx, userId)
}

func (r *userRepository) UpdatePassword(ctx context.Context, id, hashedPassword string) error {
	return r.queries.UpdatePassword(ctx, db.UpdatePasswordParams{
		ID:       id,
		Password: hashedPassword,
	})
}
