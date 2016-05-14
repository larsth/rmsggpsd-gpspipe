package cache

import (
	"testing"
	"time"

	"github.com/larsth/go-gpsfix"
	"github.com/larsth/go-rmsggpsbinmsg"
)

const float64TruncatePrecision = 4

func TestMkBinMsg(t *testing.T) {
	var (
		altitude  = float32(1.0)
		latitude  = float32(2.0)
		longitude = float32(3.0)
		fixMode   = gpsfix.Fix3D
		date      time.Time
		m         *binmsg.Message
	)
	date = time.Date(2016, time.May, 5, 18, 9, 5, 0, time.Local)
	m = mkBinMsg(altitude, latitude, longitude, fixMode, date)

	if !date.Equal(m.TimeStamp.Time) {
		t.Errorf("Got: m.Timestap.Time='%s'. Want: date='%s'",
			m.TimeStamp.Time.String(), date.String())
	}
	if m.Gps.FixMode != fixMode {
		t.Errorf("Got: '%s'. Want: '%s'",
			m.Gps.FixMode.String(), fixMode.String())
	}
	if m.Gps.Altitude != altitude {
		t.Errorf("Altitude; Got: float32(%d). Want: float32(%d)",
			m.Gps.Altitude, altitude)
	}
	if m.Gps.Latitude != latitude {
		t.Errorf("Latitude; Got: float32(%d). Want: float32(%d)",
			m.Gps.Latitude, latitude)
	}
	if m.Gps.Longitude != longitude {
		t.Errorf("Longitude; Got: float32(%d). Want: float32(%d)",
			m.Gps.Longitude, longitude)
	}
}

func TestMkFixNotSeenMessage(t *testing.T) {
	var (
		m = MkFixNotSeenMessage()
	)

	if m.TimeStamp.Time.IsZero() {
		t.Error("Got: Timestap.Time is zero. ",
			"Want: a non zero time.Time")
	}
	if m.Gps.FixMode != gpsfix.FixNotSeen {
		t.Errorf("Got: %s. Want: FixNotSeen",
			m.Gps.FixMode.String())
	}
	if m.Gps.Altitude != float32(0.0) {
		t.Errorf("Got altitude: float32(%d). %s",
			m.Gps.Altitude,
			"Want altitude: float32(0.0)")
	}
	if m.Gps.Latitude != float32(0.0) {
		t.Errorf("Got: latitude float32(%d). %s",
			m.Gps.Latitude,
			"Want latitude: float32(0.0)")
	}
	if m.Gps.Longitude != float32(0.0) {
		t.Errorf("Got longitude: float32(%d). %s",
			m.Gps.Longitude,
			"Want longitude: float32(0.0)")
	}
}

var (
	binMsgMessage1 = &binmsg.Message{
		TimeStamp: binmsg.TimeStamp{
			Time: time.Date(2015, time.May, 6, 19, 8, 4, 0, time.Local),
		},
		Gps: binmsg.Gps{
			FixMode:   gpsfix.Fix2D,
			Latitude:  float32(3.0),
			Longitude: float32(1.0),
			Altitude:  float32(4.0),
		},
	}

	binMsgMessage2 = &binmsg.Message{
		TimeStamp: binmsg.TimeStamp{
			Time: time.Date(2016, time.May, 5, 18, 9, 5, 0, time.Local),
		},
		Gps: binmsg.Gps{
			FixMode:   gpsfix.Fix3D,
			Latitude:  float32(1.0),
			Longitude: float32(2.0),
			Altitude:  float32(3.0),
		},
	}
)

func TestBinMsgPut1(t *testing.T) {
	var (
		cached            *BinMsg
		cachedMessage     *binmsg.Message = nil
		inputMessage      *binmsg.Message = nil
		wantCachedMessage *binmsg.Message = nil
		wantOk                            = false
		gotOk             bool
		s                 []string
		ok                bool
		a                 string
	)
	cached = new(BinMsg)
	if s, ok = binmsg.IsEqual(cached.m, cachedMessage, float64TruncatePrecision); !ok {
		for _, a = range s {
			t.Error(a)
		}
	}

	gotOk = cached.Put(inputMessage)
	if gotOk != wantOk {
		t.Errorf("\tGot: %v\n\tWant %v", gotOk, wantOk)
	}
	if s, ok = binmsg.IsEqual(cached.m, wantCachedMessage, float64TruncatePrecision); !ok {
		for _, a = range s {
			t.Error(a)
		}
	}
}

func TestBinMsgPut2(t *testing.T) {
	var (
		cached            *BinMsg
		cachedMessage     *binmsg.Message = nil
		inputMessage      *binmsg.Message = binMsgMessage1
		wantCachedMessage *binmsg.Message = binMsgMessage1
		wantOk                            = true
		gotOk             bool
		s                 []string
		ok                bool
		a                 string
	)
	cached = new(BinMsg)
	cached.m = cachedMessage

	gotOk = cached.Put(inputMessage)
	if gotOk != wantOk {
		t.Errorf("\tGot: %v\n\tWant %v", gotOk, wantOk)
	}
	if s, ok = binmsg.IsEqual(cached.m, wantCachedMessage, float64TruncatePrecision); !ok {
		for _, a = range s {
			t.Error(a)
		}
	}
}

func TestBinMsgPut3(t *testing.T) {
	var (
		cached            *BinMsg
		cachedMessage     *binmsg.Message = binMsgMessage1
		inputMessage      *binmsg.Message = binMsgMessage2
		wantCachedMessage *binmsg.Message = binMsgMessage2
		wantOk                            = true
		gotOk             bool
		s                 []string
		ok                bool
		a                 string
	)
	cached = new(BinMsg)
	cached.m = cachedMessage

	gotOk = cached.Put(inputMessage)
	if gotOk != wantOk {
		t.Errorf("\tGot: %v\n\tWant %v", gotOk, wantOk)
	}
	if s, ok = binmsg.IsEqual(cached.m, wantCachedMessage, float64TruncatePrecision); !ok {
		for _, a = range s {
			t.Error(a)
		}
	}
}

func TestBinMsgPut4(t *testing.T) {
	var (
		cached            *BinMsg
		cachedMessage     *binmsg.Message = binMsgMessage2
		inputMessage      *binmsg.Message = binMsgMessage1
		wantCachedMessage *binmsg.Message = binMsgMessage2
		wantOk                            = false
		gotOk             bool
		s                 []string
		ok                bool
		a                 string
	)
	cached = new(BinMsg)
	cached.m = cachedMessage

	gotOk = cached.Put(inputMessage)
	if gotOk != wantOk {
		t.Errorf("\tGot: %v\n\tWant %v", gotOk, wantOk)
	}
	if s, ok = binmsg.IsEqual(cached.m, wantCachedMessage, float64TruncatePrecision); !ok {
		for _, a = range s {
			t.Error(a)
		}
	}
}

func TestBinMsgGet1(t *testing.T) {
	var (
		cached     = new(BinMsg)
		gotMessage = binMsgMessage1
		s          []string
		ok         bool
		a          string
	)
	cached.m = binMsgMessage1
	gotMessage = cached.Get()
	if s, ok = binmsg.IsEqual(cached.m, gotMessage, float64TruncatePrecision); !ok {
		for _, a = range s {
			t.Error(a)
		}
	}
}

func TestBinMsgGet2(t *testing.T) {
	var (
		cached                        = new(BinMsg)
		cachedMessage *binmsg.Message = nil
		gotMessage    *binmsg.Message = nil
		s             []string
		ok            bool
		a             string
	)
	cached.m = cachedMessage
	gotMessage = cached.Get()
	if s, ok = binmsg.IsEqual(cachedMessage, gotMessage, float64TruncatePrecision); !ok {
		for _, a = range s {
			t.Error(a)
		}
	}
}
