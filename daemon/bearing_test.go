package daemon

import (
	"math"
	"strconv"
	"testing"
	"time"

	"github.com/larsth/go-gpsfix"
	"github.com/larsth/go-rmsggpsbinmsg"
)

type tCalcBearing struct {
	Latitude1   float32
	Longitude1  float32
	Latitude2   float32
	Longitude2  float32
	WantBaering float64
	WantErr     error
}

var tableCalcBearing = []*tCalcBearing{
	//[0] Fra Nordpolen til Nordpolen
	&tCalcBearing{
		Latitude1:   0.0,
		Longitude1:  0.0,
		Latitude2:   0.0,
		Longitude2:  0.0,
		WantBaering: math.NaN(),
		WantErr:     ErrSameLocation,
	},
	//[1] From   HAB                (lat1 lon1)=(55.69147 12.61800)
	//    to     Malta University   (lat2 lon2)=(35.90142 14.48474)
	&tCalcBearing{
		Latitude1:   55.69147,
		Longitude1:  12.61800,
		Latitude2:   35.90142,
		Longitude2:  14.48574,
		WantBaering: -1.921492, // 249.9066 degrees
		WantErr:     nil,
	},
	//[2] From   Malta University   (lat1 lon1)=(35.90142 14.48474)
	//    to     HAB                (lat2 lon2)=(55.69147 12.61800)
	&tCalcBearing{
		Latitude1:   35.90142,
		Longitude1:  14.48474,
		Latitude2:   55.69147,
		Longitude2:  12.618000,
		WantBaering: -3.1161, // 181.4582 degrees
		WantErr:     nil,
	},
}

func TestCalcBearing(t *testing.T) {
	const altitude = float32(0.0)
	var (
		u                     time.Time
		this                  *binmsg.Message
		other                 *binmsg.Message
		gotBearing            float64
		gotBearingStr         string
		gotBearingDegrees     float64
		gotBearingDegreesStr  string
		wantBearingStr        string
		wantBearingDegrees    float64
		wantBearingDegreesStr string
		gotErr                error
		ok                    bool
		s                     string
	)
	u = time.Date(2016, 5, 20, 13, 32, 0, 0, time.UTC)
	for i, td := range tableCalcBearing {
		this = binmsg.MkBinMsg(altitude, td.Latitude1, td.Longitude1,
			gpsfix.Fix3D, u)
		other = binmsg.MkBinMsg(altitude, td.Latitude2, td.Longitude2,
			gpsfix.Fix3D, u)

		wantBearingStr = strconv.FormatFloat(td.WantBaering, 'f', 4, 32)
		wantBearingDegrees = ((td.WantBaering * 180) / math.Pi) + 360.0
		wantBearingDegreesStr = strconv.FormatFloat(wantBearingDegrees, 'f', 4, 32)

		gotBearing, gotErr = calcBearing(this, other)
		gotBearingStr = strconv.FormatFloat(gotBearing, 'f', 4, 32)
		gotBearingDegrees = ((gotBearing * 180) / math.Pi) + 360.0
		gotBearingDegreesStr = strconv.FormatFloat(gotBearingDegrees, 'f', 4, 32)
		t.Logf("INFO [%d]: Want bearing:'%s'(want bearing degrees: '%s'), `Got bearing: '%s' (got bearing degrees: '%s')",
			i, wantBearingStr, wantBearingDegreesStr, gotBearingStr, gotBearingDegreesStr)

		s, ok = errorTestI(gotErr, td.WantErr, i, "tableCalcBearing")
		if !ok {
			t.Error(s)
		}
		s, ok = isSameFloat64TestI(td.WantBaering, gotBearing, i, "tableCalcBearing")
		if !ok {
			t.Errorf("Bearing test failed: %s", s)
		}
	}
}

//func TestUpdateBearingCache(t *testing.T) {}

//func TestBearingGoRoutine(t *testing.T) {}
