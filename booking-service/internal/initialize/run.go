package initialize

import (
	"github.com/gin-gonic/gin"
	"github.com/leminhthai/train-ticket/booking-service/global"
	"go.uber.org/zap"
)

func Run() *gin.Engine {
	LoadConfig()
	InitLogger()

	global.Logger.Info("Init Logger Successfully!", zap.String("ok", "success"))

	InitMySql()
	InitRedis()

	app := InitServiceApp()
	InitHandlers(app)
	InitWorker(app)

	r := InitRouter()

	return r
}
