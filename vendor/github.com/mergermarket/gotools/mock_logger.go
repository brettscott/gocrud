package tools

import (
	"errors"
	"fmt"
)

// MockLogger provides a basic mock of the logger object.
type MockLogger struct {
	calls []LoggerCall
}

// Call is a single call to a logger method. It has the method name and the arguments it was called with
type LoggerCall struct {
	Method string
	Args   LoggerArgs
}

// Args are the list of arguments to a single logger method
type LoggerArgs struct {
	Msg string
}

// Info is a mock info method
func (ml *MockLogger) Info(args ...interface{}) {
	ml.call("Info", args)
}

// Info is a mock info method
func (ml *MockLogger) Error(args ...interface{}) {
	ml.call("Error", args)
}

// Debug is a mock debug method
func (ml *MockLogger) Debug(args ...interface{}) {
	ml.call("Debug", args)
}

func (ml *MockLogger) Call() (c LoggerCall, err error) {
	if len(ml.calls) == 0 {
		return c, errors.New("No calls made")
	}
	return ml.calls[0], nil
}

func (ml *MockLogger) LastCall() *LoggerCall {
	if len(ml.calls) == 0 {
		return nil
	}
	return &ml.calls[len(ml.calls)-1]
}

func (ml *MockLogger) call(method string, args ...interface{}) {
	msg := fmt.Sprint(args...)
	ml.calls = append(ml.calls, LoggerCall{method, LoggerArgs{msg}})
}
