package api

// StatsDer is the interface for the DataDog StatsD methods
type StatsDer interface {
	Histogram(name string, value float64, tags ...string)
	Gauge(name string, value float64, tags ...string)
	Incr(name string, tags ...string)
}
