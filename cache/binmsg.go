package cache

import (
	"sync"
	"time"

	"github.com/larsth/go-gpsfix"
	"github.com/larsth/go-rmsggpsbinmsg"
)

type BinMsg struct {
	mutex sync.RWMutex
	m     *binmsg.Message
}

func mkBinMsg(altitude, latitude, longitude float32,
	fixMode gpsfix.FixMode, t time.Time) *binmsg.Message {
	var (
		m = new(binmsg.Message)
	)

	m.TimeStamp.Time = t
	m.Gps.Altitude = altitude
	m.Gps.Latitude = latitude
	m.Gps.Longitude = longitude
	m.Gps.FixMode = fixMode

	return m
}

func MkFixNotSeenMessage() *binmsg.Message {
	return mkBinMsg(float32(0.0), float32(0.0), float32(0.0),
		gpsfix.FixNotSeen, time.Now().UTC())
}

func (b *BinMsg) Put(m *binmsg.Message) (ok bool) {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	if m == nil {
		return false
	}

	if m != nil && b.m == nil {
		//Fast caching
		b.m = m
		return true
	}

	if m.TimeStamp.Time.After(b.m.TimeStamp.Time) {
		//Cache message 'm'
		b.m = m
		return true
	}

	//Forget old message 'm'
	return false
}

func (b *BinMsg) Get() *binmsg.Message {
	b.mutex.RLock()
	defer b.mutex.RUnlock()

	return b.m
}
