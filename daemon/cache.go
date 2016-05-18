package daemon

import (
	"math"
	"time"

	"github.com/larsth/go-rmsggpsbinmsg"
	"github.com/larsth/rmsggpsd-gpspipe/cache"
)

var (
	thisGpsCache  *cache.BinMsg
	otherGpsCache *cache.BinMsg
	bearingCache  *cache.Bearing
	thisChan      chan *binmsg.Message
	otherChan     chan *binmsg.Message
)

func init() {
	thisChan = make(chan *binmsg.Message, 1)
	otherChan = make(chan *binmsg.Message, 1)

	thisGpsCache, _ = cache.NewBinMsg(thisChan)
	otherGpsCache, _ = cache.NewBinMsg(otherChan)

	bearingCache = new(cache.Bearing)
	bearingCache.Put(math.NaN(), time.Unix(0, 0))
}
