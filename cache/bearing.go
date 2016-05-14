package cache

import (
	"sync"
	"time"
)

type Bearing struct {
	mutex   sync.RWMutex
	bearing float64
	t       time.Time
}

func (b *Bearing) Put(bearing float64, t time.Time) bool {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	if b.t.Before(t) {
		b.bearing = bearing
		b.t = t
		return true
	}
	return false
}

func (b *Bearing) Get() (float64, time.Time) {
	b.mutex.RLock()
	defer b.mutex.RUnlock()

	return b.bearing, b.t
}
