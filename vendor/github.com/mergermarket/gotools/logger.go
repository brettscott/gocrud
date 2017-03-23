package tools

import (
	"os"
	"strings"

	"github.com/Sirupsen/logrus"
	"runtime"
)

type Logger interface {
	Info(msg ...interface{})
	Error(msg ...interface{})
	Debug(msg ...interface{})
	Warn(msg ...interface{})
}

// Logger logs messages in a structured format in prod and pretty colours in local.
type logrusLogger struct {
	log    *logrus.Logger
	fields logrus.Fields
}

// Info should be used to log key application events.
func (l *logrusLogger) Info(msg ...interface{}) {
	l.log.WithFields(withFileAndLine(l.fields)).Info(msg)
}

// Error should be used to log events that need to be actioned on immediately.
func (l *logrusLogger) Error(msg ...interface{}) {
	l.log.WithFields(withFileAndLine(l.fields)).Error(msg)
}

// Debug can be used to log events for local development.
func (l *logrusLogger) Debug(msg ...interface{}) {
	l.log.WithFields(withFileAndLine(l.fields)).Debug(msg)
}

// Warn is for when something bad happened but doesnt need instant action.
func (l *logrusLogger) Warn(msg ...interface{}) {
	l.log.WithFields(withFileAndLine(l.fields)).Warn(msg)
}

func withFileAndLine(fields logrus.Fields) logrus.Fields {
	newFields := make(map[string]interface{})
	_, file, line, ok := runtime.Caller(2)
	for k, v := range fields {
		newFields[k] = v
	}
	if ok {
		newFields["file"] = file
		newFields["line"] = line
	} else {
		newFields["file"] = "Unknown"
		newFields["line"] = "Unknown"
	}
	return newFields
}

// NewLogger returns a new structured logger.
func NewLogger(isLocal bool) Logger {
	logger := logrus.New()

	if strings.ToLower(os.Getenv("LOG_LEVEL")) == "debug" {
		logger.Level = logrus.DebugLevel
	}

	if !isLocal {
		logger.Formatter = &logrus.JSONFormatter{}
	}

	return &logrusLogger{
		log: logger,
		fields: logrus.Fields{
			"component": getComponentName(),
			"env":       getEnv(),
		},
	}
}

func getComponentName() string {
	if name := os.Getenv("COMPONENT_NAME"); len(name) > 0 {
		return name
	}
	return "a-service-has-no-name"
}

func getEnv() string {
	if env := os.Getenv("ENV_NAME"); len(env) > 0 {
		return env
	}
	return "local"
}
