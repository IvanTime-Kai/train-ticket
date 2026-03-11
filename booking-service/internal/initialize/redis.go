package initialize

import (
	"context"
	"fmt"

	"github.com/leminhthai/train-ticket/booking-service/global"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

var ctx = context.Background()

func InitRedis() {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", global.Config.Redis.Host, global.Config.Redis.Port),
		Password: global.Config.Redis.Password,
		DB:       global.Config.Redis.DB,
	})

	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		global.Logger.Error("Init Redis Failed!!!", zap.Error(err))
		panic(err)
	}

	global.Logger.Info("Init Redis Successfully!!!")
	global.Rdb = rdb
}
