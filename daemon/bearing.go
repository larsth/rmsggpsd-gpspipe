package daemon

import (
	"math"

	"github.com/larsth/rmsggpsd-gpspipe/errors"
)

type LatLon struct {
	Lat float64
	Lon float64
}

/*
bearing is a function that calculates the _inital_ bearing
from point p1(lat1, lon1) to point p2(lat2, lon2) by using the
haversine algorithm.

The bearing changes at any point between point 1 to
point 2 and vise versa.

The returned float64 value is the initial bearing in degrees (not radians).
*/
func bearing(p1, p2 LatLon) (float64, error) {
	var (
		// dLat    float64 : dLat not used. Why?
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

	//dLat = p2.Lat - p1.Lat : dLat not used. Why?
	dLon = p2.Lon - p1.Lon
	cosLat1 = math.Cos(p1.Lat)
	cosLat2 = math.Cos(p2.Lat)
	sinLat1 = math.Sin(p1.Lat)
	sinLat2 = math.Sin(p2.Lat)

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
		Note that dlat is NEVER used

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
