package cache

import (
	"testing"
	"time"

	"github.com/larsth/go-gpsfix"
	"github.com/larsth/go-rmsggpsbinmsg"
)

const float64TruncatePrecision = 4

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
		cached            *BinMsg         = new(BinMsg)
		inputMessage      *binmsg.Message = binMsgMessage1
		wantCachedMessage *binmsg.Message = binMsgMessage1
		wantOk                            = true
		gotOk             bool
		s                 []string
		ok                bool
		a                 string
	)
	cached.m = nil //the cached *binmsg.Message is nil

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
