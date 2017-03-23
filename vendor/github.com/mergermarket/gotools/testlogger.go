package tools

type T interface {
	Log(args ...interface{})
}

// TestLogger accepts the testing package so you wont be bombarded with logs
// when your tests pass but if they fail you will see what's going on.
type TestLogger struct {
	T T
}

// Info logs info to the test logger.
func (testLogger TestLogger) Info(msg ...interface{}) {
	testLogger.T.Log("[Info]", msg)
}

// Debug logs debug to the test logger.
func (testLogger TestLogger) Debug(msg ...interface{}) {
	testLogger.T.Log("[Debug]", msg)
}

// Error logs error to the test logger.
func (testLogger TestLogger) Error(msg ...interface{}) {
	testLogger.T.Log("[Error]", msg)
}

// Warn logs warn to the test logger.
func (testLogger TestLogger) Warn(msg ...interface{}) {
	testLogger.T.Log("[Warn]", msg)
}
