package daemon

import (
	"math"
	"time"

	"github.com/larsth/rmsggpsd-gpspipe/cache"
)

var (
	thisGpsCache  *cache.BinMsg
	otherGpsCache *cache.BinMsg
	bearingCache  *cache.Bearing
)

func init() {
	thisGpsCache = new(cache.BinMsg)
	thisGpsCache.Put(cache.MkFixNotSeenMessage())

	otherGpsCache = new(cache.BinMsg)
	otherGpsCache.Put(cache.MkFixNotSeenMessage())

	bearingCache = new(cache.Bearing)
	bearingCache.Put(float64(4*math.Pi), time.Unix(0, 0))
}
