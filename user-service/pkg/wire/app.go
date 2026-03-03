package wire

import (
	"github.com/leminhthai/train-ticket/user-service/internal/service"
)

type App struct {
	UserService service.UserService
}
