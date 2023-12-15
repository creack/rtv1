package math2

import (
	"fmt"
	"math"
)

// Point represents a 2D point.
type Point struct{ X, Y float64 }

// Pt is a helper to create a point.
func Pt[T ~int | ~float32 | ~float64](x, y T) Point {
	return Point{X: float64(x), Y: float64(y)}
}

// String representation.
func (p Point) String() string {
	p.X = math.Round(p.X*100) / 100
	p.Y = math.Round(p.Y*100) / 100
	return fmt.Sprintf("(%v,%v)", p.X, p.Y)
}

// Add p2 to p.
func (p Point) Add(p2 Point) Point {
	return Point{
		X: p.X + p2.X,
		Y: p.Y + p2.Y,
	}
}

// Sub p2 from p.
func (p Point) Sub(p2 Point) Point {
	return Point{
		X: p.X - p2.X,
		Y: p.Y - p2.Y,
	}
}

// Mul p and p2.
func (p Point) Mul(p2 Point) Point {
	return Point{
		X: p.X * p2.X,
		Y: p.Y * p2.Y,
	}
}

// Div p with p2.
// NOTE: Go handles div by 0 by setting the value to Inf.
func (p Point) Div(p2 Point) Point {
	return Point{
		X: p.X / p2.X,
		Y: p.Y / p2.Y,
	}
}

// Scale the point by the given factor.
func (p Point) Scale(n float64) Point {
	return Point{
		X: p.X * n,
		Y: p.Y * n,
	}
}

// Magnitude returns the length of the line p->p2.
func (p Point) Magnitude(p2 Point) float64 { return p.Sub(p2).Norm() }

// Norm returns the point's norm. aka length aka magnitude aka hypot.
func (p Point) Norm() float64 { return math.Hypot(p.X, p.Y) }

// Rotate the point by angle.
// See math2.Rotate to rotate around a given point.
func (p Point) Rotate(alpha Angle) Point {
	// Lookup cos(alpha) and sin(alpha).
	sin, cos := math.Sincos(alpha.Radians())

	// Rotate.
	return Point{
		X: p.X*cos - p.Y*sin,
		Y: p.X*sin + p.Y*cos,
	}
}
