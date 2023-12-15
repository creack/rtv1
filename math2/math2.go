// Package math2 provides 2D math helpers.
package math2

import (
	"math"
)

// Scale is a generic helper wrapping p.Scale(n).
func Scale[T ~int | ~float32 | ~float64](p Point, n T) Point {
	return p.Scale(float64(n))
}

// GetAngle returns the angle between 0,0 -> p and 0,0 -> p2.
func (p Point) GetAngle(p2 Point) Angle {
	theta1 := math.Atan2(p.Y, p.X)
	theta2 := math.Atan2(p2.Y, p2.X)
	return NewRadAngle(theta2 - theta1)
}

// GetAngle returns the angle between lines origin->p1 and origin->p2.
func GetAngle(origin, p1, p2 Point) Angle {
	// Remove the origin offset.
	p1 = p1.Sub(origin)
	p2 = p2.Sub(origin)
	// Get the angle.
	return p1.GetAngle(p2)
}

// CoordinatesFromAngleDist returns the coordinates from 0,0
// of the point from an angle and length.
func (p Point) CoordinatesFromAngleDist(a Angle, length float64) Point {
	// Create a vector of size 'length' at the origin.
	pv := Point{length, 0}
	// Get the angle between origin and pp1.
	origAngle := Point{}.GetAngle(p)
	// Rotate the vector with the origRad angle to be at the origin->p1 angle.
	pv1 := pv.Rotate(origAngle)
	// Then rotate the vector with the relative angle to be at the origin->p2 angle.
	pv2 := pv1.Rotate(a)

	// Add the origin offset to the result.
	return pv2
}

// CoordinatesFromAngleDist returns the coordinates from 0,0
// of a point p1 offseted by origin from an angle and length.
func CoordinatesFromAngleDist(origin, p1 Point, a Angle, length float64) Point {
	return p1.
		Sub(origin).                         // Remove the origin offset.
		CoordinatesFromAngleDist(a, length). // Get the coordinates.
		Add(origin)                          // Add the origin offset back.
}

// Rotate the point around origin by angle.
func Rotate(p Point, alpha Angle, origin Point) Point {
	// Remove the origin offset.
	p = p.Sub(origin)
	// Rotate.
	p = p.Rotate(alpha)
	// Add back the origin offset.
	return p.Add(origin)
}
