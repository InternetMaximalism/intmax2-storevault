package logger

import (
	"intmax2-store-vault/internal/logger"
	"intmax2-store-vault/internal/logger/logrus"
)

func New(logLevel, timeFormat string, logJSON, logLines bool) logger.Logger {
	return logrus.New(logLevel, timeFormat, logJSON, logLines)
}
