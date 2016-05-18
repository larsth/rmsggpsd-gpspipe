package cache

import (
	"errors"
	"sync"
	"time"

	"github.com/larsth/go-gpsfix"
	"github.com/larsth/go-rmsggpsbinmsg"
)

var ErrFuncIsNil = errors.New("Function 'f func()' is nil")

type BinMsg struct {
	mutex sync.RWMutex
	m     *binmsg.Message
	c     chan *binmsg.Message
}

func NewBinMsg(c chan *binmsg.Message) (*BinMsg, error) {
	b := new(BinMsg)
	b.m = MkFixNotSeenMessage()
	b.c = c
	return b, nil
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
		b.c <- m
		return true
	}

	if m.TimeStamp.Time.After(b.m.TimeStamp.Time) {
		//Cache message 'm'
		b.m = m
		b.c <- m
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
