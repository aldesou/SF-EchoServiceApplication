package logger

type Logger interface {
	Close() error
	LogInfo(correlationId string, message string)
	LogWarning(correlationId string, message string)
	LogError(correlationId string, message string)
}
