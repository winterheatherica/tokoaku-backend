package cache

import "time"

const (
	TickInterval1h  = 1 * time.Hour
	TickInterval3h  = 3 * time.Hour
	TickInterval6h  = 6 * time.Hour
	TickInterval12h = 12 * time.Hour
	TickInterval24h = 24 * time.Hour

	SleepOnError = 30 * time.Second
)
