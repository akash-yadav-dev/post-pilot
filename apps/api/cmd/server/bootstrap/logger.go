package bootstrap

import "post-pilot/packages/logger"

type Logger = logger.Logger

func NewLogger() (Logger, error) {
	return logger.NewFromEnv()
}
