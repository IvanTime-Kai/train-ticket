package initialize

import (
	"github.com/leminhthai/train-ticket/booking-service/global"
	"github.com/leminhthai/train-ticket/booking-service/pkg/logger"
)

func InitLogger() {
	global.Logger = logger.NewLogger(global.Config.Logger)
}
