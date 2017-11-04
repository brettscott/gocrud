package crud

import "testing"

// Logger interface for logging
type Logger interface {
	Info(msg ...interface{})
	Error(msg ...interface{})
	Debug(msg ...interface{})
	Warn(msg ...interface{})
}

type T interface {
	Log(args ...interface{})
}

// NewTestLog to create a logger for use during tests
func NewTestLog(t *testing.T) TestLog {
	return TestLog{
		T: t,
	}
}

// TestLog accepts the testing package so you wont be bombarded with logs
// when your tests pass but if they fail you will see what's going on.
type TestLog struct {
	T T
}

// Info logs info to the test logger.
func (testLog TestLog) Info(msg ...interface{}) {
	testLog.T.Log("[Info]", msg)
}

// Debug logs debug to the test logger.
func (testLog TestLog) Debug(msg ...interface{}) {
	testLog.T.Log("[Debug]", msg)
}

// Error logs error to the test logger.
func (testLog TestLog) Error(msg ...interface{}) {
	testLog.T.Log("[Error]", msg)
}

// Warn logs warn to the test logger.
func (testLog TestLog) Warn(msg ...interface{}) {
	testLog.T.Log("[Warn]", msg)
}
