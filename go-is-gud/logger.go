package main

import (
	"go.uber.org/zap"
)

func NewLogger() *zap.SugaredLogger {
	config := zap.NewProductionConfig()
	config.OutputPaths = []string{"stdout"}
	logger, err := config.Build()
	if err != nil {
		panic(err)
	}

	return logger.Sugar()
}
