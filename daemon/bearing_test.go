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
		WantBaering: -1.921492, //-110.0934 grader ???
		WantErr:     nil,
	},
	//[2] Fra Nordpolen til Nordpolen
	//	&tCalcBearing{
	//		Latitude1:   0.0,
	//		Longitude1:  0.0,
	//		Latitude2:   0.0,
	//		Longitude2:  0.0,
	//		WantBaering: math.NaN(),
	//		WantErr:     ErrSameLocation,
	//	},
}

func TestCalcBearing(t *testing.T) {
	const altitude = float32(0.0)
	var (
		u          time.Time
		this       *binmsg.Message
		other      *binmsg.Message
		gotBearing float64
		gotErr     error
		ok         bool
		s          string
	)
	u = time.Date(2016, 5, 20, 13, 32, 0, 0, time.UTC)
	for i, td := range tableCalcBearing {
		this = binmsg.MkBinMsg(altitude, td.Latitude1, td.Longitude1,
			gpsfix.Fix3D, u)
		other = binmsg.MkBinMsg(altitude, td.Latitude2, td.Longitude2,
			gpsfix.Fix3D, u)

		gotBearing, gotErr = calcBearing(this, other)

		s, ok = errorTestI(gotErr, td.WantErr, i, "tableCalcBearing")
		if !ok {
			t.Error(s)
		}
		s, ok = isSameFloat64TestI(td.WantBaering, gotBearing, i, "tableCalcBearing")
		if !ok {
			gbd := (gotBearing * 180) / math.Pi
			gotBearingDeg := strconv.FormatFloat(gbd, 'f', 4, 32)
			wbd := (td.WantBaering * 180) / math.Pi
			wantBaeringDeg := strconv.FormatFloat(wbd, 'f', 4, 32)
			t.Errorf("Bearing test failed: %s; (Bearing degrees, got=%s, want=%s)",
				s, gotBearingDeg, wantBaeringDeg)
		}
	}
}

//func TestUpdateBearingCache(t *testing.T) {}

//func TestBearingGoRoutine(t *testing.T) {}
