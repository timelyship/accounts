package application

import (
	"fmt"
	"go.uber.org/zap"
)

func NewLogger(traceID, spanID string) *zap.Logger {
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
	)
}

func NewTraceableLogger(l interface{}, exists bool) *zap.Logger {
	if exists {
		return l.(*zap.Logger)
	}
	return NewLogger("<?>", "<?>")
}
