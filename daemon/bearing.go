package daemon

import (
	"log"
	"math"
	"time"

	"github.com/larsth/go-gpsfix"
	"github.com/larsth/go-rmsggpsbinmsg"
	"github.com/larsth/rmsggpsd-gpspipe/errors"
)

/*
calcBearing is a function that calculates the _inital_ bearing
from point p1(lat1, lon1) to point p2(lat2, lon2) by using the
haversine algorithm.

The bearing changes at any point between point 1 to
point 2 and vise versa.

The returned float64 value is the initial bearing in radians (not degrees).
*/
func calcBearing(this, other *binmsg.Message) (float64, error) {
	/*
		tc1=mod(atan2(sin(lon2-lon1)*cos(lat2),
		   	cos(lat1)*sin(lat2)-sin(lat1)*cos(lat2)*cos(lon2-lon1)), 2*pi)

		Algorithm is from:
		    http://mathforum.org/library/drmath/view/55417.html : Post #2

		The algorithm had been translated to this pseudo code:
		   	dlon = lon2 -lon1
		   	sin_dlon = sin(dlon)
		   	cos_lat2 = cos(lat2)
		   	cos_lat1 = cos(lat1)
		   	sin_lat2 = sin(lat2)
		   	sin_lat1 = sin(lat1)
		   	cos_dlon = cos(dlon)
		   	atan2_x = sin_dlon*cos_lat2
		   	atan2_y1 = cos_lat1*sin_lat2
		   	atan2_y2 = sin_lat1*cos_lat2*cos_dlon
		   	atan2_y = atan2_y1 - atan2_y2
		    mod_x = atan2(atan2_x, atan2_y)
		    mod_y = 2*pi
		   	bearing=mod(mod_x, mod_y)
	*/
	var (
		lon1, lat1, lon2, lat2           float64
		dLon                             float64
		sinLat2, sinLat1, sinDlon        float64
		cosLat2, cosLat1, cosDlon        float64
		atan2x, atan2y1, atan2y2, atan2y float64
		modX, modY, bearing              float64
		err                              error
	)

	lat1 = this.Gps.Lat()
	lat2 = other.Gps.Lat()
	lon1 = this.Gps.Lon()
	lon2 = other.Gps.Lon()

	dLon = lon2 - lon1

	sinDlon = math.Sin(dLon)
	cosLat2 = math.Cos(lat2)
	cosLat1 = math.Cos(lat1)
	sinLat2 = math.Sin(lat2)
	sinLat1 = math.Sin(lat1)
	cosDlon = math.Cos(dLon)

	atan2x = sinDlon * cosLat2
	atan2y1 = cosLat1 * sinLat2
	atan2y2 = sinLat1 * cosLat2 * cosDlon
	atan2y = atan2y1 - atan2y2

	if modX, err = atan2(atan2x, atan2y); err != nil {
		return math.NaN(), errors.Trace(err)
	}
	modY = math.Pi * 2

	bearing = math.Mod(modX, modY)

	return bearing, nil
}

func updateBearingCache(this, other *binmsg.Message, logger *log.Logger) {
	var (
		useTime    time.Time
		a, b, c, d bool
		bearing    float64
		err        error
	)

	if this.TimeStamp.Time.Before(other.TimeStamp.Time) {
		useTime = this.TimeStamp.Time
	} else {
		useTime = other.TimeStamp.Time
	}

	a = this.Gps.FixMode == gpsfix.FixNotSeen
	b = this.Gps.FixMode == gpsfix.FixNone
	c = this.Gps.FixMode == gpsfix.FixNotSeen
	d = this.Gps.FixMode == gpsfix.FixNone

	if a || b || c || d {
		bearing = math.NaN()
	} else {
		//FIXME
		/* Also use sentinel values?, ie. : special float64 values outside 360
		   degrees (2*Pi rad) to give extra information to http request functions?
		*/
		bearing, err = calcBearing(this, other)
		if err != nil {
			logger.Println(err)
		}
	}
	bearingCache.Put(bearing, useTime)
}

func bearingGoRoutine(logger *log.Logger) {
	var (
		this  = binmsg.MkFixNotSeenMessage()
		other = binmsg.MkFixNotSeenMessage()
	)

	updateBearingCache(this, other, logger)
	for {
		select {
		case this = <-thisChan:
			updateBearingCache(this, other, logger)
		case other = <-otherChan:
			updateBearingCache(this, other, logger)
		}
	}
}
