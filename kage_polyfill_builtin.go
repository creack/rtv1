package main

import "math"

// This file is a wrapper for Kage. It mirror the shaderlib_builtin.kage file to allow
// for the shader to be compile in Go.
// In this part, we have the default built-in functions and types.

type vec2 struct{ x, y float }

func (v vec2) uniform() []float32 {
	return []float32{float32(v.x), float32(v.y)}
}

type vec3 struct {
	vec2
	z float
}

func (v vec3) uniform() []float32 {
	xy := v.vec2.uniform()
	return []float32{xy[0], xy[1], float32(v.z)}
}

type vec4 struct {
	vec3
	w float
}

func (v vec4) uniform() []float32 {
	xyz := v.vec3.uniform()
	return []float32{xyz[0], xyz[1], xyz[2], float32(v.w)}
}

type float = float64

type mat4 [4]vec4

func length3(v vec3) float {
	return sqrt(dot3(v, v))
}

func normalize3(v vec3) vec3 {
	mag := length3(v)
	if mag > 0 {
		invLen := 1.0 / mag
		return newVec3(
			v.x*invLen,
			v.y*invLen,
			v.z*invLen,
		)
	}
	return v
}

func dot3(v1, v2 vec3) float {
	return v1.x*v2.x + v1.y*v2.y + v1.z*v2.z
}

func cos(in float) float    { return math.Cos(in) }
func sin(in float) float    { return math.Sin(in) }
func tan(in float) float    { return math.Tan(in) }
func acos(in float) float   { return math.Acos(in) }
func asin(in float) float   { return math.Asin(in) }
func atan(in float) float   { return math.Atan(in) }
func sqrt(in float) float   { return math.Sqrt(in) }
func pow(in, n float) float { return math.Pow(in, n) }
func floor(in float) float  { return math.Floor(in) }

const pi = math.Pi

var (
	_ = pi
	_ = cos(0)
	_ = sin(0)
	_ = tan(0)
	_ = acos(0)
	_ = asin(0)
	_ = atan(0)
	_ = sqrt(0)
	_ = pow(0, 0)
	_ = floor(0)
)

func newVec3(x, y, z float) vec3 {
	var v vec3
	v.x = x
	v.y = y
	v.z = z
	return v
}

func newVec4(x, y, z, w float) vec4 {
	var v vec4
	v.x = x
	v.y = y
	v.z = z
	v.w = w
	return v
}

func newMat4(a, b, c, d vec4) mat4 {
	return mat4{a, b, c, d}
}

func sub3(v1, v2 vec3) vec3 {
	return newVec3(
		v1.x-v2.x,
		v1.y-v2.y,
		v1.z-v2.z,
	)
}

func add3(v1, v2 vec3) vec3 {
	return newVec3(
		v1.x+v2.x,
		v1.y+v2.y,
		v1.z+v2.z,
	)
}

func cross3(v1, v2 vec3) vec3 {
	return newVec3(
		v1.y*v2.z-v1.z*v2.y,
		v1.z*v2.x-v1.x*v2.z,
		v1.x*v2.y-v1.y*v2.x,
	)
}

func scale3(v1 vec3, f float) vec3 {
	return newVec3(
		v1.x*f,
		v1.y*f,
		v1.z*f,
	)
}

func mul4(v1, v2 vec4) vec4 {
	return newVec4(
		v1.x*v2.x,
		v1.y*v2.y,
		v1.z*v2.z,
		1, //v1.w*v2.w,
	)
}

func scale4(v1 vec4, f float) vec4 {
	return newVec4(
		v1.x*f,
		v1.y*f,
		v1.z*f,
		1, //v1.w*f,
	)
}

func add4(v1, v2 vec4) vec4 {
	return newVec4(
		v1.x+v2.x,
		v1.y+v2.y,
		v1.z+v2.z,
		1, //v1.w+v2.w,
	)
}
