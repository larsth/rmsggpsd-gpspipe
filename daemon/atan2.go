package daemon

import (
	"math"

	"github.com/larsth/rmsggpsd-gpspipe/errors"
)

var ErrSameLocation = errors.New("The 2 coordinates is at the same location")

//atan2 is the arctangent 2 algoritm
/* See https://en.wikipedia.org/wiki/Atan2
Quote:
In a variety of computer languages, the function atan2 is the arctangent
function with two arguments. The purpose of using two arguments instead of one
is to gather information on the signs of the inputs in order to return the
appropriate quadrant of the computed angle, which is not possible for the
single-argument arctangent function. It also avoids the problems of division by
zero.

For any real number (e.g., floating point) arguments x and y not both equal to
zero, atan2(y, x) is the angle in radians between the positive x-axis of a
plane and the point given by the coordinates (x, y) on it.
The angle is positive for counter-clockwise angles (upper half-plane, y > 0),
and negative for clockwise angles (lower half-plane, y < 0).
*/
func atan2(x, y float64) (float64, error) {
	if y > 0.0 {
		if x > 0.0 {
			return math.Atan((y / x)), nil
		}
		if x < 0.0 {
			return (180.0 - math.Atan((-y / x))), nil
		}
		if x == 0.0 {
			return 90.0, nil
		}
		return math.NaN(), errors.Errorf(
			"X is not <, > or equal to 0.0. x:=\"%d\"\n\n", x)
	}
	if y < 0.0 {
		if x > 0.0 {
			return (-1 * math.Atan((-y / x))), nil
		}
		if x < 0.0 {
			return (math.Atan((y / x)) - 180.0), nil
		}
		if x == 0.0 {
			return 270.0, nil
		}
		return math.NaN(), errors.Errorf(
			"X is not <, > or equal to 0.0. x:=\"%d\"\n\n", y)
	}
	if y == 0.0 {
		if x > 0.0 {
			return 0.0, nil
		}
		if x < 0.0 {
			return 180.0, nil
		}
		if x == 0.0 {
			return math.NaN(), ErrSameLocation
		}
		return math.NaN(), errors.Errorf(
			"X is not <, > or equal to 0.0. x:=\"%d\"\n\n", y)
	}
	// y is not <, >, or equal to 0.0!
	return math.NaN(), errors.Errorf(
		"Y is not <, > or equal to 0.0. y:=\"%d\"\n\n", y)
}
