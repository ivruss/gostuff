package logger

import (
	"fmt"
	"go.uber.org/zap"
)

func ConstructLogger(env string) (*zap.Logger, error) {
	var logger *zap.Logger
	var err error

	switch env {
	case "dev":
		logger, err = zap.NewDevelopment()
		if err != nil {
			return nil, fmt.Errorf("error creating development logger: %w", err)
		}
	case "prod":
		logger, err = zap.NewProduction()
		if err != nil {
			return nil, fmt.Errorf("error creating production logger: %w", err)
		}
	default:
		return nil, fmt.Errorf("unknown environment %q", cfg.Logger.Env)
	}

	defer logger.Sync()
	return logger, nil
}
