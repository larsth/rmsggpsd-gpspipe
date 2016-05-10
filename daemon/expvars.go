package daemon

import (
	"time"

	"github.com/larsth/go-rmsggpsbinmsg"
	"github.com/larsth/rmsggpsd/expvar"
)

var (
	cacheVars               *expvar.Map
	cacheVarLat             *expvar.Float
	cacheVarLon             *expvar.Float
	cacheVarAlt             *expvar.Float
	cacheVarFixMode         *expvar.String
	cacheVarCurrentTime     *expvar.String
	cacheVarLastCompareTime *expvar.String
)

func initExpVars() {
	cacheVars = expvar.NewMap("cache")

	cacheVarLat = expvar.NewFloat("lat")
	cacheVarLat.Set(fixNotSeenMessage.Gps.Lat())

	cacheVarLon = expvar.NewFloat("lon")
	cacheVarLon.Set(fixNotSeenMessage.Gps.Lon())

	cacheVarAlt = expvar.NewFloat("alt")
	cacheVarLat.Set(fixNotSeenMessage.Gps.Alt())

	cacheVarFixMode = expvar.NewString("fixmode")
	cacheVarFixMode.Set(fixNotSeenMessage.Gps.FixMode.String())

	cacheVarCurrentTime = expvar.NewString("current-time")
	cacheVarCurrentTime.Set(fixNotSeenMessage.TimeStamp.Time.Format(time.RFC3339))

	cacheVarLastCompareTime = expvar.NewString("last-compare-time")
	cacheVarLastCompareTime.Set(fixNotSeenMessage.TimeStamp.Time.Format(time.RFC3339))

	cacheVars.Set("lat", cacheVarLat)
	cacheVars.Set("lon", cacheVarLon)
	cacheVars.Set("alt", cacheVarAlt)
	cacheVars.Set("fixmode", cacheVarFixMode)
	cacheVars.Set("current-time", cacheVarCurrentTime)
	cacheVars.Set("last-compare-time", cacheVarLastCompareTime)
}

func setExpVarCached(m *binmsg.Message, lastCompareTime time.Time) {
	cacheVarLat.Set(m.Gps.Lat())
	cacheVarLon.Set(m.Gps.Lon())
	cacheVarAlt.Set(m.Gps.Alt())
	cacheVarFixMode.Set(m.Gps.FixMode.String())

	cacheVarCurrentTime.Set(m.TimeStamp.Time.Format(time.RFC3339))
	cacheVarLastCompareTime.Set(lastCompareTime.Format(time.RFC3339))
}

func setExpVarLastCompareTime(lastCompareTime time.Time) {
	cacheVarLastCompareTime.Set(lastCompareTime.Format(time.RFC3339))
}
