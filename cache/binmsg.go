package cache

import (
	"errors"
	"sync"

	"github.com/larsth/go-rmsggpsbinmsg"
)

var ErrFuncIsNil = errors.New("Function 'f func()' is nil")

type BinMsg struct {
	mutex sync.RWMutex
	m     *binmsg.Message
	C     chan *binmsg.Message
}

func NewBinMsg(c chan *binmsg.Message) (*BinMsg, error) {
	b := new(BinMsg)
	b.m = binmsg.MkFixNotSeenMessage()
	b.C = c
	return b, nil
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
		if b.C != nil {
			b.C <- m
		}
		return true
	}

	if m.TimeStamp.Time.After(b.m.TimeStamp.Time) {
		//Cache message 'm'
		b.m = m

		if b.C != nil {
			b.C <- m
		}
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
