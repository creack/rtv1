package math2

import (
	"fmt"
	"math"
)

// Angle helper type.
type Angle float64

// NewRadAngle returns a new angle from a radians value.
func NewRadAngle[T ~int | ~float32 | ~float64](radians T) Angle {
	return Angle(radians).Normalize()
}

// NewDegAngle returns a new angle from a degrees values.
func NewDegAngle[T ~int | ~float32 | ~float64](degrees T) Angle {
	return NewRadAngle(float64(degrees) * math.Pi / 180)
}

// String in degrees.
func (a Angle) String() string { return fmt.Sprint(math.Round(a.Degrees()*100) / 100) }

// Radians returns the radians value of the angle.
func (a Angle) Radians() float64 { return float64(a) }

// Degrees returns the degress value of the angle.
func (a Angle) Degrees() float64 { return float64(a) * 180. / math.Pi }

// Normalize the angle on -pi to +pi.
// i.e. 270 degrees -> -90 degrees.
func (a Angle) Normalize() Angle {
	// Module 2pi (360). Doesn't change the sign.
	a = Angle(math.Mod(a.Radians(), 2*math.Pi))

	if a > math.Pi {
		return a - 2*math.Pi
	}
	if a < -math.Pi {
		return a + 2*math.Pi
	}
	return a
}

// Between checks if the angle is in between the parameter range.
func (a Angle) Between(start, end Angle) bool {
	// Normalize angles.
	a = a.Normalize()
	start = start.Normalize()
	end = end.Normalize()

	// Check if angle is between start and end.
	if start <= end {
		return a >= start && a <= end
	}

	// Handle the case where the range spans 0 degrees.
	return a >= start || a <= end
}
