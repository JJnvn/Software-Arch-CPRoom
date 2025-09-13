package logs

import (
	"fmt"
	"os"

	"go.uber.org/zap"
)

const (
	DEV  = "development"
	PROD = "production"
)

func NewLogger() (*zap.Logger, error) {
	return newLoggerFactory(os.Getenv("SERVER_ENV"))
}

func newLoggerFactory(env string) (*zap.Logger, error) {
	var logger *zap.Logger
	var err error

	switch env {
	case DEV:
		logger, err = zap.NewDevelopment()
		if err != nil {
			return nil, err
		}
	case PROD:
		logger, err = zap.NewProduction()
		if err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("invalid environment: %s", env)
	}

	return logger, nil
}
