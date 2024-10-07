package utils

import (
	"os"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	LOG_LEVEL string = "info"
)

func CreateLogger() *zap.Logger {
	logLevel, logLevelSet := os.LookupEnv("LOG_LEVEL")
	if logLevelSet {
		LOG_LEVEL = strings.ToLower(logLevel)
	}

	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.TimeKey = "timestamp"
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder

	zapLogLevel := zap.InfoLevel

	switch LOG_LEVEL {
	case "debug":
		zapLogLevel = zap.DebugLevel
	case "info":
		zapLogLevel = zap.InfoLevel
	case "warn":
		zapLogLevel = zap.WarnLevel
	case "error":
		zapLogLevel = zap.ErrorLevel
	case "panic":
		zapLogLevel = zap.DPanicLevel
	case "fatal":
		zapLogLevel = zap.FatalLevel
	}

	config := zap.Config{
		Level:             zap.NewAtomicLevelAt(zapLogLevel),
		Development:       false,
		DisableCaller:     false,
		DisableStacktrace: false,
		Sampling:          nil,
		Encoding:          "json",
		EncoderConfig:     encoderCfg,
		OutputPaths: []string{
			"stderr",
		},
		ErrorOutputPaths: []string{
			"stderr",
		},
		InitialFields: map[string]interface{}{
			"pid": os.Getpid(),
		},
	}

	return zap.Must(config.Build())
}
