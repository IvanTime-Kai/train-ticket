package initialize

import (
	"github.com/leminhthai/train-ticket/train-service/global"
	"github.com/leminhthai/train-ticket/train-service/pkg/logger"
)

func InitLogger() {
	global.Logger = logger.NewLogger(global.Config.Logger)
}
