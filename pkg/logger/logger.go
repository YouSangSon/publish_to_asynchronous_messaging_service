package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	Sugar *zap.SugaredLogger
)

// InitLogger initializes the logger
func InitLogger(loggerConfig *zap.Config, mode string) error {
	var logger *zap.Logger
	var err error

	if mode == "dev" {
		loggerConfig.EncoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
		logger, err = loggerConfig.Build(zap.AddCaller())
		if err != nil {
			return err
		}
	} else {
		logger, err = loggerConfig.Build()
		if err != nil {
			return err
		}
	}

	Sugar = logger.Sugar()

	if err != nil {
		return err
	}

	return nil
}
