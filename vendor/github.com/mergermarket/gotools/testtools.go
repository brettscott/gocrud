package tools

func NewTestTools(t T) (Logger, StatsD) {
	logger := TestLogger{T: t}
	statsd, _ := NewStatsD(NewStatsDConfig(false, logger))
	return logger, statsd
}
