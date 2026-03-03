package initialize

import (
	"github.com/leminhthai/train-ticket/user-service/global"
	"github.com/leminhthai/train-ticket/user-service/internal/handler"
	"github.com/leminhthai/train-ticket/user-service/pkg/wire"
)

func InitServices() {
	initUserServices()
}

func initUserServices() {
	app := wire.InitializeApp(global.Mdb)
	handler.InitAuth(app.UserService)
	handler.InitUser(app.UserService)
}
