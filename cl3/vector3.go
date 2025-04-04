package main

import (
	"math"
)

// Vector3 represents a 3D vector or point
type Vector3 struct {
	X, Y, Z float64
}

// Add returns the sum of two vectors
func (v Vector3) Add(other Vector3) Vector3 {
	return Vector3{v.X + other.X, v.Y + other.Y, v.Z + other.Z}
}

// Sub returns the difference between two vectors
func (v Vector3) Sub(other Vector3) Vector3 {
	return Vector3{v.X - other.X, v.Y - other.Y, v.Z - other.Z}
}

// Mul returns the product of a vector and a scalar
func (v Vector3) Mul(scalar float64) Vector3 {
	return Vector3{v.X * scalar, v.Y * scalar, v.Z * scalar}
}

// Div returns the division of a vector by a scalar
func (v Vector3) Div(scalar float64) Vector3 {
	return Vector3{v.X / scalar, v.Y / scalar, v.Z / scalar}
}

// Dot returns the dot product of two vectors
func (v Vector3) Dot(other Vector3) float64 {
	return v.X*other.X + v.Y*other.Y + v.Z*other.Z
}

// Cross returns the cross product of two vectors
func (v Vector3) Cross(other Vector3) Vector3 {
	return Vector3{
		v.Y*other.Z - v.Z*other.Y,
		v.Z*other.X - v.X*other.Z,
		v.X*other.Y - v.Y*other.X,
	}
}

// Length returns the magnitude of the vector
func (v Vector3) Length() float64 {
	return math.Sqrt(v.Dot(v))
}

// Normalize returns the unit vector
func (v Vector3) Normalize() Vector3 {
	length := v.Length()
	if length == 0 {
		return Vector3{0, 0, 0}
	}
	return v.Div(length)
}

// Reflect returns the reflection vector
func (v Vector3) Reflect(normal Vector3) Vector3 {
	return v.Sub(normal.Mul(2 * v.Dot(normal)))
}

// RotateX rotates a point around the X axis
func RotateX(point Vector3, angle float64) Vector3 {
	cosA := math.Cos(angle)
	sinA := math.Sin(angle)

	return Vector3{
		X: point.X,
		Y: point.Y*cosA - point.Z*sinA,
		Z: point.Y*sinA + point.Z*cosA,
	}
}

// RotateY rotates a point around the Y axis
func RotateY(point Vector3, angle float64) Vector3 {
	cosA := math.Cos(angle)
	sinA := math.Sin(angle)

	return Vector3{
		X: point.X*cosA + point.Z*sinA,
		Y: point.Y,
		Z: -point.X*sinA + point.Z*cosA,
	}
}
