package daemon

import (
	"sync"
	"time"

	"github.com/larsth/go-gpsfix"
	"github.com/larsth/go-rmsggpsbinmsg"
	"github.com/larsth/rmsggpsd/errors"
)

type binMsgCache struct {
	mutex sync.RWMutex
	m     *binmsg.Message
}

func (c *binMsgCache) Get() (*binmsg.Message, error) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	if c.m == nil {
		return fixNotSeenMessage, errors.New("Cached message is nil")
	}
	return c.m, nil

}

func (c *binMsgCache) isAfter(t time.Time) bool {
	if c.m == nil {
		return false
	}

	if t.After(c.m.TimeStamp.Time) {
		//If t.Time is after the c.m.TimeStamp.Time, then:
		return true
	}
	//else:
	return false
}

func (c *binMsgCache) Put(m *binmsg.Message) (ok bool) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if c.isAfter(m.TimeStamp.Time) {
		//Cache message 'm'
		c.m = m
		setExpVarCached(m, time.Now().UTC())
		ok = true

	} else {
		//Forget message 'm'
		setExpVarLastCompareTime(time.Now().UTC())
		ok = false
	}
	return
}

var (
	cache             *binMsgCache
	fixNotSeenMessage *binmsg.Message
)

func initFixNotSeenMessage() {
	cache = new(binMsgCache)

	fixNotSeenMessage = new(binmsg.Message)
	cache.m = fixNotSeenMessage

	fixNotSeenMessage.TimeStamp.Time = time.Now().UTC()
	fixNotSeenMessage.Gps.FixMode = gpsfix.FixNotSeen
	fixNotSeenMessage.Gps.Altitude = float32(0.0)
	fixNotSeenMessage.Gps.Latitude = float32(0.0)
	fixNotSeenMessage.Gps.Longitude = float32(0.0)
}
