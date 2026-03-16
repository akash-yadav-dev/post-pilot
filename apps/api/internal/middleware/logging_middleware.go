package middleware

type LoggingMiddleware struct {
	// You can add fields here if needed, e.g., a logger instance
}

func NewLoggingMiddleware() *LoggingMiddleware {
	return &LoggingMiddleware{}
}
