package main

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func initLogger(logLevel string, isLogJson bool) *zap.Logger {
	opts := zap.NewProductionConfig()
	opts.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	if !isLogJson {
		opts.Encoding = "console"
		opts.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}

	lvl := zap.InfoLevel
	lvlErr := lvl.UnmarshalText([]byte(logLevel))
	opts.Level = zap.NewAtomicLevelAt(lvl)

	logger, err := opts.Build()
	if err != nil {
		panic(fmt.Sprintf("Failed to initialize logger: %s\n", err))
	}

	if lvlErr != nil {
		logger.Error("Invalid log level. Expected debug, info, warn, error, dpanic, panic, fatal", zap.Error(lvlErr))
	}

	return logger
}
