package cache

import (
	"strconv"
	"sync"
	"time"
)

type Bearing struct {
	mutex      sync.RWMutex
	bearing    float64
	t          time.Time
	bearingStr string
	tStr       string
}

func (b *Bearing) Put(bearing float64, t time.Time) bool {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	if b.t.Before(t) {
		b.bearing = bearing
		b.t = t
		b.bearingStr = strconv.FormatFloat(bearing, 'f', -1, 32)
		b.tStr = t.Format(time.RFC3339)
		return true
	}
	return false
}

func (b *Bearing) Get() (bearing float64, t time.Time, bearingStr string, tStr string) {
	b.mutex.RLock()
	defer b.mutex.RUnlock()

	return b.bearing, b.t, b.bearingStr, tStr
}
