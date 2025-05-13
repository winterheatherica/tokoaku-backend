package persistent

import (
	"sync"
)

var (
	persistentRedisPrefix string
	persistentOnce        sync.Once
)
