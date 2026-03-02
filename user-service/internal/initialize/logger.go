package initialize

import (
	"github.com/leminhthai/train-ticket/user-service/global"
	"github.com/leminhthai/train-ticket/user-service/pkg/logger"
)

func InitLogger() {
	global.Logger = logger.NewLogger(global.Config.Logger)
}
