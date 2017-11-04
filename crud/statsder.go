package crud

import "testing"

// StatsDer is the interface for the DataDog StatsD methods
type StatsDer interface {
	Histogram(name string, value float64, tags ...string)
	Gauge(name string, value float64, tags ...string)
	Incr(name string, tags ...string)
}

// NewTestStatsD to create a logger for use during tests
func NewTestStatsD(t *testing.T) TestStatsD {
	return TestStatsD{
		T: t,
	}
}

// TestStatsD accepts the testing package so you wont be bombarded with logs
// when your tests pass but if they fail you will see what's going on.
type TestStatsD struct {
	T T
}

// Histogram logs
func (testStatsD TestStatsD) Histogram(name string, value float64, tags ...string) {
	testStatsD.T.Log("[Histogram]", name, value, tags)
}

// Gauge logs
func (testStatsD TestStatsD) Gauge(name string, value float64, tags ...string) {
	testStatsD.T.Log("[Gauge]", name, value, tags)
}

// Incr logs
func (testStatsD TestStatsD) Incr(name string, tags ...string) {
	testStatsD.T.Log("[Incr]", name, tags)
}

