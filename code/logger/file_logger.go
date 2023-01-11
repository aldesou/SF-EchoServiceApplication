package logger

import (
	"fmt"
	"time"

	"gopkg.in/natefinch/lumberjack.v2"
)

type fileLoggerImpl struct {
	logger *lumberjack.Logger
}

func NewFileLogger(logPath string) Logger {
	logger := lumberjack.Logger{
		Filename:   logPath,
		MaxSize:    1, // MB
		MaxBackups: 3,
		MaxAge:     30, //days
		Compress:   false,
	}
	return &fileLoggerImpl{
		logger: &logger,
	}
}

func (logger *fileLoggerImpl) Close() error {
	return logger.logger.Close()
}

func (logger *fileLoggerImpl) LogInfo(correlationId string, message string) {
	fmt.Fprintf(logger.logger, "[%s][INFO][%s] %s\n", time.Now().UTC().String(), correlationId, message)
}

func (logger *fileLoggerImpl) LogWarning(correlationId string, message string) {
	fmt.Fprintf(logger.logger, "[%s][WARN][%s] %s\n", time.Now().UTC().String(), correlationId, message)
}

func (logger *fileLoggerImpl) LogError(correlationId string, message string) {
	fmt.Fprintf(logger.logger, "[%s][ERROR][%s] %s\n", time.Now().UTC().String(), correlationId, message)
}
