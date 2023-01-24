package logging

import (
	"context"
	"log"

	"go.uber.org/zap"
)

type loggerKeyType int
const loggerKey loggerKeyType = iota

var logger *zap.Logger

func init() {
	l, err := zap.NewDevelopment()
	if err != nil {
		log.Fatal(err)
	}
	logger = l
}

func WithContext(ctx context.Context) *zap.Logger {
	if ctx == nil {
		return logger
	}

	if ctxLogger, ok := ctx.Value(loggerKey).(zap.Logger); ok {
		return &ctxLogger
	}
	return logger
}
