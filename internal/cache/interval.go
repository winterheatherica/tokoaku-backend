package cache

import "time"

const (
	TickInterval3m  = 3 * time.Minute
	TickInterval5m  = 5 * time.Minute
	TickInterval10m = 10 * time.Minute
	TickInterval15m = 15 * time.Minute
	TickInterval30m = 30 * time.Minute

	TickInterval1h  = 1 * time.Hour
	TickInterval3h  = 3 * time.Hour
	TickInterval6h  = 6 * time.Hour
	TickInterval12h = 12 * time.Hour
	TickInterval24h = 24 * time.Hour

	SleepOnError = 30 * time.Second
)
