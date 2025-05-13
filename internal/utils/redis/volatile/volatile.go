package volatile

import (
	"sync"
)

var (
	volatileRedisPrefix string
	volatileOnce        sync.Once
)
