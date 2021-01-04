package application

import (
	"go.uber.org/zap"
)

var (
	Logger = NewLogger()
)

func NewLogger() *zap.Logger {
	config := zap.NewProductionConfig()
	config.OutputPaths = []string{
		"stdout",
	}

	config.EncoderConfig.LevelKey = "level"
	config.EncoderConfig.TimeKey = "timestamp"
	config.EncoderConfig.CallerKey = "caller"
	config.EncoderConfig.MessageKey = "message"

	l, _ := config.Build()
	return l
}

func SyncLogger() {
	_ = Logger.Sync()
}
