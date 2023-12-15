// Package math3 provides basic 3D math utils.
package math3

import (
	"math"
)

// Vec is a 3d vector.
type Vec struct {
	X, Y, Z float64
}

// ScaleAll the vector by the given factor.
func (v Vec) ScaleAll(scale float64) Vec {
	return Vec{
		X: v.X * scale,
		Y: v.Y * scale,
		Z: v.Z * scale,
	}
}

// ScaleZ scales only Z.
func (v Vec) ScaleZ(scale float64) Vec {
	return Vec{
		X: v.X,
		Y: v.Y,
		Z: v.Z * scale,
	}
}

// Translate the vector.
func (v Vec) Translate(offset Vec) Vec {
	return Vec{
		X: v.X + offset.X,
		Y: v.Y + offset.Y,
		Z: v.Z + offset.Z,
	}
}

// Rotate the vector.
func (v Vec) Rotate(angle Vec) Vec {
	v = v.MultiplyMatrix(GetRotationMatrix(angle.Z, AxisZ))
	v = v.MultiplyMatrix(GetRotationMatrix(angle.X, AxisX))
	v = v.MultiplyMatrix(GetRotationMatrix(angle.Y, AxisY))
	return v
}

// MultiplyMatrix multiplies the given matrix with the current vector.
func (v Vec) MultiplyMatrix(m Matrix) Vec {
	return Vec{
		X: v.X*m[0][0] + v.Y*m[0][1] + v.Z*m[0][2],
		Y: v.X*m[1][0] + v.Y*m[1][1] + v.Z*m[1][2],
		Z: v.X*m[2][0] + v.Y*m[2][1] + v.Z*m[2][2],
	}
}

// Matrix is a 3x3 matrix.
type Matrix [3][3]float64

// Multiply 2 3x3 matrices.
func (m Matrix) Multiply(m2 Matrix) Matrix {
	var result Matrix

	// Multiplying matrices and storing result.
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			for k := 0; k < 3; k++ {
				result[i][j] += m[i][k] * m2[k][j]
			}
		}
	}
	return result
}

// Axis enum type.
type Axis byte

// Axis enum values.
const (
	AxisNone Axis = iota
	AxisX
	AxisY
	AxisZ
)

// GetRotationMatrix returns the populated matrix to rotate along the given axis.
//
// Ref: https://en.wikipedia.org/wiki/Rotation_matrix#Basic_3D_rotations.
func GetRotationMatrix(deg float64, axis Axis) Matrix {
	c := math.Cos(deg)
	s := math.Sin(deg)

	//nolint:exhaustive // False positive.
	switch axis {
	case AxisX:
		return Matrix{
			{1, 0, 0},
			{0, c, -s},
			{0, s, c},
		}
	case AxisY:
		return Matrix{
			{c, -s, 0},
			{0, 1, 0},
			{-s, 0, c},
		}
	case AxisZ:
		return Matrix{
			{c, -s, 0},
			{s, c, 0},
			{0, 0, 1},
		}
	default:
		panic("invalid axis")
	}
}
