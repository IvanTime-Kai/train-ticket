package global

import (
	"database/sql"

	"github.com/leminhthai/train-ticket/user-service/pkg/logger"
	"github.com/leminhthai/train-ticket/user-service/pkg/setting"
	"github.com/redis/go-redis/v9"
)

var (
	Config setting.Settings
	Logger *logger.LoggerZap
	Mdb    *sql.DB
	Rdb    *redis.Client
)
