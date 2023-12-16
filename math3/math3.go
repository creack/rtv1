// Package math3 provides basic 3D math utils.
package math3

import (
	"math"
)

// Point is an alias for Vector for clarity.
type Point = Vector

// Vector is a 3d vector.
type Vector struct {
	X, Y, Z float64
}

// Vec is a helper to create a vector.
func Vec[T ~int | ~float64](x, y, z T) Vector {
	return Vector{X: float64(x), Y: float64(y), Z: float64(z)}
}

// Norma normalizes.
func (v Vector) Norm() Vector {
	mag := v.Magnitude()
	div := -1.
	if mag != 0 {
		div = 1 / mag
	}
	return v.ScaleAll(div)
}

// Magnitude is the length.
func (v Vector) Magnitude() float64 {
	return math.Sqrt(v.Dot(v))
}

// ScaleAll the vector by the given factor.
func (v Vector) ScaleAll(scale float64) Vector {
	return Vector{
		X: v.X * scale,
		Y: v.Y * scale,
		Z: v.Z * scale,
	}
}

// ScaleZ scales only Z.
func (v Vector) ScaleZ(scale float64) Vector {
	return Vector{
		X: v.X,
		Y: v.Y,
		Z: v.Z * scale,
	}
}

// Translate the vector.
func (v Vector) Translate(offset Vector) Vector {
	return Vector{
		X: v.X + offset.X,
		Y: v.Y + offset.Y,
		Z: v.Z + offset.Z,
	}
}

// Add v2 to v.
func (v Vector) Add(v2 Vector) Vector { return v.Translate(v2) }

// Sub v2 from v.
func (v Vector) Sub(v2 Vector) Vector { return v.Translate(Vector{X: -v2.X, Y: -v2.Y, Z: -v2.Z}) }

// Mul v2 and v.
func (v Vector) Mul(v2 Vector) Vector {
	return v.Translate(Vector{X: v.X * v2.X, Y: v.Y * v2.Y, Z: v.Z * v2.Z})
}

// Dot v2 and v.
func (v Vector) Dot(ov Vector) float64 {
	return v.X*ov.X + v.Y*ov.Y + v.Z*ov.Z
}

// Cross v2 and v.
func (v Vector) Cross(v2 Vector) Vector {
	return Vec(
		v.Y*v2.Z-v.Z*v2.Y,
		v.Z*v2.X-v.X*v2.Z,
		v.X*v2.Y-v.Y*v2.X,
	)
}

// Rotate the vector.
func (v Vector) Rotate(angle Vector) Vector {
	v = v.MultiplyMatrix(GetRotationMatrix(angle.Z, AxisZ))
	v = v.MultiplyMatrix(GetRotationMatrix(angle.X, AxisX))
	v = v.MultiplyMatrix(GetRotationMatrix(angle.Y, AxisY))
	return v
}

// MultiplyMatrix multiplies the given matrix with the current vector.
func (v Vector) MultiplyMatrix(m Matrix) Vector {
	return Vector{
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
