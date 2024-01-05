package logger

import (
	"errors"
	"os"

	"go.uber.org/zap"
)

var (
	ErrMissingEnvVar = errors.New("missing environment variable")
)

func New() (*zap.Logger, error) {
	var logger *zap.Logger
	var err error

	environment := os.Getenv("ENVIRONMENT")
	if environment == "" {
		return nil, ErrMissingEnvVar
	}

	switch environment {
	case "development":
		logger, err = zap.NewDevelopment()
	case "production":
		logger, err = zap.NewProduction()
	default:
		logger = zap.NewExample()
		logger.Warn("Unsupported environment", zap.String("environment", environment), zap.Strings("supportedEnvironments", []string{"development", "production"}), zap.String("action", "creating example zap logger"))
	}

	return logger, err
}
