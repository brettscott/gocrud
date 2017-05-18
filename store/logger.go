package store

// Logger interface for logging
type Logger interface {
	Info(msg ...interface{})
	Error(msg ...interface{})
	Debug(msg ...interface{})
	Warn(msg ...interface{})
}
