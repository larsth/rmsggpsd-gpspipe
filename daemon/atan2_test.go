package daemon

import (
	"math"
	"testing"

	"github.com/juju/errors"

	//"github.com/larsth/rmsggpsd-gpspipe/errors"
)

type atan2TestType struct {
	FuncMkWantZ func(td *atan2TestType)
	X           float64
	Y           float64
	WantZ       float64
	WantErr     error
}

var atan2TestTable = []*atan2TestType{
	// [0], y > 0, x > 0 :
	&atan2TestType{
		FuncMkWantZ: func(td *atan2TestType) {
			td.WantZ = math.Atan((td.X / td.Y))
		},
		X:       math.Pi,
		Y:       math.Pi,
		WantErr: nil,
	},
	// [1], y > 0, x < 0 :
	&atan2TestType{
		FuncMkWantZ: func(td *atan2TestType) {
			td.WantZ = (math.Pi - math.Atan((-td.Y / td.X)))
		},
		X:       -1.0 * math.Pi,
		Y:       math.Pi,
		WantErr: nil,
	},
	// [2], y > 0, x == 0 :
	&atan2TestType{
		X:       float64(0.0),
		Y:       math.Pi,
		WantZ:   (math.Pi / 2),
		WantErr: nil,
	},
	// [3], y > 0, x == NaN :
	&atan2TestType{
		FuncMkWantZ: func(td *atan2TestType) {
			td.WantErr = errors.Errorf(
				"X is not <, > or equal to 0.0. x:=\"%d\"\n\n", td.X)
		},
		WantZ: math.NaN(),
		X:     math.NaN(),
		Y:     math.Pi,
	},
	// [4], y < 0, x > 0 :
	&atan2TestType{
		FuncMkWantZ: func(td *atan2TestType) {
			td.WantZ = (-1 * math.Atan((-td.Y / td.X)))
		},
		X:       math.Pi,
		Y:       -1.0 * math.Pi,
		WantErr: nil,
	},
	// [5], y < 0, x < 0 :
	&atan2TestType{
		FuncMkWantZ: func(td *atan2TestType) {
			td.WantZ = (math.Atan((td.X / td.Y)) - math.Pi)
		},
		X:       -1.0 * math.Pi,
		Y:       -1.0 * math.Pi,
		WantErr: nil,
	},
	// [6], y < 0, x == 0 :
	&atan2TestType{
		X:       0.0,
		Y:       -1.0 * math.Pi,
		WantZ:   ((3 * math.Pi) / 2),
		WantErr: nil,
	},
	// [7], y < 0, x == NaN :
	&atan2TestType{
		FuncMkWantZ: func(td *atan2TestType) {
			td.WantErr = errors.Errorf(
				"X is not <, > or equal to 0.0. x:=\"%d\"\n\n", td.X)
		},
		X:       math.NaN(),
		Y:       -1.0 * math.Pi,
		WantZ:   math.NaN(),
		WantErr: nil,
	},
	// [8], y == 0, x > 0 :
	&atan2TestType{
		X:       math.Pi,
		Y:       0.0,
		WantZ:   0.0,
		WantErr: nil,
	},
	// [9], y == 0, x < 0 :
	&atan2TestType{
		X:       -1.0 * math.Pi,
		Y:       0.0,
		WantZ:   math.Pi,
		WantErr: nil,
	},
	// [10], y == 0, x == 0 :
	&atan2TestType{
		X:       0.0,
		Y:       0.0,
		WantZ:   math.NaN(),
		WantErr: ErrSameLocation,
	},
	// [11], y == 0, x == NaN :
	&atan2TestType{
		FuncMkWantZ: func(td *atan2TestType) {
			td.WantErr = errors.Errorf("%s. x:=\"%v\"",
				"X is not <, > or equal to 0.0", td.X)
		},
		X:     math.NaN(),
		Y:     0.0,
		WantZ: math.NaN(),
	},
	// [12], y == ?, x == ? :
	&atan2TestType{
		FuncMkWantZ: func(td *atan2TestType) {
			td.WantErr = errors.Errorf("%s  x:=\"%v\", y:=\"%v\"",
				"X,Y is not <, > or equal to 0.0.", td.X, td.Y)
		},
		X:     math.NaN(),
		Y:     math.NaN(),
		WantZ: math.NaN(),
	},
}

func TestAtan2(t *testing.T) {
	var (
		gotZ   float64
		gotErr error
		s      string
		ok     bool
	)
	for i, td := range atan2TestTable {
		if td.FuncMkWantZ != nil {
			td.FuncMkWantZ(td)
		}
		gotZ, gotErr = atan2(td.X, td.Y)

		s, ok = errorTestI(gotErr, td.WantErr, i, "atan2TestTable")
		if !ok {
			t.Error(s)
		}

		s, ok = isSameFloat64TestI(td.WantZ, gotZ, i, "atan2TestTable")
		if !ok {
			t.Error(s)
		}
	}
}
