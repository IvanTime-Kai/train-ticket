package global

import (
	"database/sql"

	"github.com/leminhthai/train-ticket/train-service/pkg/logger"
	"github.com/leminhthai/train-ticket/train-service/pkg/setting"
)

var (
	Config setting.Settings
	Mdb    *sql.DB
	Logger *logger.LoggerZap
)
