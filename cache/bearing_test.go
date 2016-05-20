package cache

import (
	"testing"
	"time"

	"github.com/larsth/go-rmsggpsbinmsg"
)

func TestBearingPut1(t *testing.T) {
	var (
		b             = new(Bearing)
		cachedBearing = float64(1.0)
		newBearing    = float64(2.0)
		cachedTime    = time.Unix(int64(10000000000), 0)
		newTime       = time.Unix(int64(20000000000), 0)
		wantOk        = true
		gotOk         bool
		cTxt          string
		nTxt          string
		ok            bool
	)

	b.bearing = cachedBearing
	b.t = cachedTime

	gotOk = b.Put(newBearing, newTime)
	if gotOk != wantOk {
		t.Errorf("\tGot Ok: %v\n\tWant Ok %v", gotOk, wantOk)
	}
	cTxt, nTxt, ok = binmsg.IsSameFloat64(cachedBearing, newBearing, float64TruncatePrecision)
	if ok {
		t.Errorf("cached bearing %f (%s) is not equal to %f (%s) , precision: %d",
			cachedBearing, cTxt, newBearing, nTxt, float64TruncatePrecision)
	}
	if !b.t.Equal(newTime) {
		t.Errorf("Cached time %s is not equal to the wanted time: %s",
			b.t.String(), newTime.String())
	}
}

func TestBearingPut2(t *testing.T) {
	var (
		b             = new(Bearing)
		cachedBearing = float64(1.0)
		newBearing    = float64(2.0)
		cachedTime    = time.Unix(int64(20000000000), 0)
		newTime       = time.Unix(int64(10000000000), 0)
		wantOk        = false
		gotOk         bool
		cTxt          string
		nTxt          string
		ok            bool
	)

	b.bearing = cachedBearing
	b.t = cachedTime

	gotOk = b.Put(newBearing, newTime)
	if gotOk != wantOk {
		t.Errorf("\tGot Ok: %v\n\tWant Ok %v", gotOk, wantOk)
	}
	cTxt, nTxt, ok = binmsg.IsSameFloat64(cachedBearing, newBearing, float64TruncatePrecision)
	if ok {
		t.Errorf("Cached bearing %f (%s) is not equal to the wanted bearing: %f (%s) , precision: %d",
			cachedBearing, cTxt, newBearing, nTxt, float64TruncatePrecision)
	}
	if b.t.Equal(newTime) {
		t.Errorf("Cached time %s is not equal to the wanted time: %s",
			b.t.String(), newTime.String())
	}
}

func TestBearingGet(t *testing.T) {
	var (
		b             = new(Bearing)
		cachedBearing = float64(1.0)
		cachedTime    = time.Unix(int64(20000000000), 0)
		gotBearing    float64
		gotTime       time.Time
		cTxt          string
		nTxt          string
		ok            bool
	)
	b.bearing = cachedBearing
	b.t = cachedTime
	gotBearing, gotTime, _, _ = b.Get()
	cTxt, nTxt, ok = binmsg.IsSameFloat64(cachedBearing, gotBearing, float64TruncatePrecision)
	if !ok {
		t.Errorf("Cached bearing %f (%s) is not equal to the wanted bearing:  %f (%s) , precision: %d",
			cachedBearing, cTxt, gotBearing, nTxt, float64TruncatePrecision)
	}
	if !gotTime.Equal(cachedTime) {
		t.Errorf("Cached time %s is not equal to the wanted time: %s",
			cachedTime, gotTime)
	}
}
