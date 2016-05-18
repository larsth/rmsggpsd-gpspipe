package daemon

import (
	"log"
	"math"
	"time"

	"github.com/larsth/go-gpsfix"
	"github.com/larsth/go-rmsggpsbinmsg"
	"github.com/larsth/rmsggpsd-gpspipe/cache"
	"github.com/larsth/rmsggpsd-gpspipe/errors"
)

type LatLon struct {
	Lat float64
	Lon float64
}

/*
calcBearing is a function that calculates the _inital_ bearing
from point p1(lat1, lon1) to point p2(lat2, lon2) by using the
haversine algorithm.

The bearing changes at any point between point 1 to
point 2 and vise versa.

The returned float64 value is the initial bearing in radians (not degrees).
*/
func calcBearing(this, other *binmsg.Message) (float64, error) {
	var (
		// dLat    float64 : dLat not used. Why?
		lat1 = this.Gps.Lat()
		lat2 = other.Gps.Lat()
		lon1 = this.Gps.Lon()
		lon2 = other.Gps.Lon()

		dLon    float64
		cosLat1 float64
		cosLat2 float64
		sinLat1 float64
		sinLat2 float64
		y       float64
		x       float64
		bearing float64
		err     error
	)

	//dLat = p2Lat - p1Lat : dLat not used. Why?
	dLon = lon2 - lon1
	cosLat1 = math.Cos(lat1)
	cosLat2 = math.Cos(lat2)
	sinLat1 = math.Sin(lat1)
	sinLat2 = math.Sin(lat2)

	y = math.Sin(dLon) * cosLat2
	x = (cosLat1 * sinLat2) - (sinLat1 * cosLat2 * math.Cos(dLon))
	if bearing, err = atan2(x, y); err != nil {
		return math.NaN(), errors.Annotate(err, "error using the atan2 algotithm")
	}
	return bearing, nil

	/*
		The algoritm below is from:
		http://mathforum.org/library/drmath/view/55417.html , where:
			tc1 is the inital bearing.

		    Note that dlat is NEVER used (why?)

		dlat = lat2 - lat1
		dlon = lon2 - lon1
		y = sin(lon2-lon1)*cos(lat2)
		x = cos(lat1)*sin(lat2)-sin(lat1)*cos(lat2)*cos(lon2-lon1)
		if y > 0 then
			if x > 0 then tc1 = arctan(y/x)
			if x < 0 then tc1 = 180 - arctan(-y/x)
			if x = 0 then tc1 = 90
		if y < 0 then
			if x > 0 then tc1 = -arctan(-y/x)
			if x < 0 then tc1 = arctan(y/x)-180
			if x = 0 then tc1 = 270
		if y = 0 then
			if x > 0 then tc1 = 0
			if x < 0 then tc1 = 180
			if x = 0 then [the 2 points are the same]
	*/
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
		//BUG
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
		this  = cache.MkFixNotSeenMessage()
		other = cache.MkFixNotSeenMessage()
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
