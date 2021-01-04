package application

import (
	"fmt"
	"go.uber.org/zap"
)

func NewLogger(traceID, spanID, userID string) *zap.Logger {
	config := zap.NewDevelopmentConfig()
	config.OutputPaths = []string{
		"stdout",
	}
	config.EncoderConfig.LevelKey = "level"
	config.EncoderConfig.TimeKey = "timestamp"
	config.EncoderConfig.CallerKey = "caller"
	config.EncoderConfig.MessageKey = fmt.Sprintf("message[%s]", traceID)
	l, _ := config.Build()
	return l.With(zap.String("traceID", traceID),
		zap.String("spanID", spanID),
		zap.String("userID", userID),
	)
}

func NewTraceableLogger(l interface{}) zap.Logger {
	return l.(zap.Logger)
}
