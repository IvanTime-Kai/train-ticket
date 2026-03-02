package logger

import (
	"os"

	"github.com/leminhthai/train-ticket/user-service/pkg/setting"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

type LoggerZap struct {
	*zap.Logger
}

func NewLogger(config setting.LoggerSetting) *LoggerZap {
	logLevel := config.Log_level

	// debug  -> info -> warning -> error -> fatal -> panic
	var level zapcore.Level

	switch logLevel {
	case "debug":
		level = zapcore.DebugLevel
	case "info":
		level = zapcore.InfoLevel
	case "warning":
		level = zapcore.WarnLevel
	case "error":
		level = zapcore.ErrorLevel
	case "fatal":
		level = zapcore.FatalLevel
	case "panic":
		level = zapcore.PanicLevel
	default:
		level = zapcore.InfoLevel
	}

	encodeLog := getEncoderLog()

	hook := lumberjack.Logger{
		Filename:   config.File_log_name,
		MaxSize:    config.Max_size,    // MB
		MaxAge:     config.Max_age,     // Auto remove (day)
		MaxBackups: config.Max_backups, // số lượng file backups
		Compress:   config.Compress,
	}

	core := zapcore.NewCore(
		encodeLog,
		zapcore.NewMultiWriteSyncer(
			zapcore.AddSync(os.Stdout),
			zapcore.AddSync(&hook),
		),
		level,
	)

	return &LoggerZap{
		Logger: zap.New(core, zap.AddCaller(), zap.AddStacktrace(zap.ErrorLevel)),
	}

}

func getEncoderLog() zapcore.Encoder {
	encodeConfig := zap.NewProductionEncoderConfig()

	encodeConfig.TimeKey = "time"
	encodeConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encodeConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encodeConfig.EncodeCaller = zapcore.ShortCallerEncoder

	return zapcore.NewJSONEncoder(encodeConfig)
}
