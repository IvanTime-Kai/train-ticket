//go:build wireinject
// +build wireinject

package wire

import (
	"database/sql"

	"github.com/google/wire"
	"github.com/leminhthai/train-ticket/user-service/internal/repository"
	"github.com/leminhthai/train-ticket/user-service/internal/service"
)

// UserSet nhóm provider cho user (queries -> repo -> service).
var UserSet = wire.NewSet(
	repository.NewUserRepository,
	service.NewUserService,
)

// InitializeApp tạo toàn bộ dependency. Một injector, một lần *sql.DB. Chạy `make wire` sau khi sửa Set/App.
func InitializeApp(sqlDB *sql.DB) *App {
	wire.Build(
		ProvideQueries,
		UserSet,
		wire.Struct(new(App), "*"),
	)
	return nil
}
