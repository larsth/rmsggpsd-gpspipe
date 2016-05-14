package daemon

import (
	"github.com/larsth/rmsggpsd-gpspipe/cache"
)

var (
	thisGpsCache  *cache.BinMsg
	otherGpsCache *cache.BinMsg
	bearingCache  *cache.Bearing
)

func init() {
	thisGpsCache = new(cache.BinMsg)
	otherGpsCache = new(cache.BinMsg)
	bearingCache = new(cache.Bearing)
}
