package main

import "fmt"

// This file is a wrapper for Kage. It mirror the shaderlib_rtv1.kage file to allow
// for the shader to be compile in Go.
// In this part, we have the RTv1 specific functions and types.

type ThingsT []mat4

type LightsT []mat4

// Globals used to pass the scene to the Fragment function.
// In shader mode, it the constructors get injected.
var (
	sceneObjects ThingsT
	sceneLights  LightsT
)

type sphere struct {
	Center vec3  `json:"center"`
	Radius float `json:"radius"`
	Color  vec4  `json:"color"`
}

func (s sphere) mat4() mat4 { return newSphere(s.Center, s.Radius, s.Color) }

func (s sphere) marshalConstructor() string {
	return fmt.Sprintf("newSphere(%s, %f, %s)", s.Center.marshalConstructor(), s.Radius, s.Color.marshalConstructor())
}

//nolint:unparam // Keeping for reference.
func getSphere(in mat4) (center vec3, radius, radius2 float, color vec4) {
	return in[0].vec3, in[0].w, in[1].w, in[2]
}

type plane struct {
	Center vec3  `json:"center"`
	Offset float `json:"offset"`
	Color  vec4  `json:"color"`
}

func (p plane) mat4() mat4 { return newPlane(p.Center, p.Offset, p.Color) }

func (p plane) marshalConstructor() string {
	return fmt.Sprintf("newPlane(%s, %f, %s)", p.Center.marshalConstructor(), p.Offset, p.Color.marshalConstructor())
}

func getPlane(in mat4) (center vec3, offset, roughness float, color vec4) {
	return in[0].vec3, in[0].w, in[1].x, in[2]
}

type light struct {
	Origin vec3 `json:"origin"`
	Color  vec4 `json:"color"`
}

func (l light) mat4() mat4 { return newLight(l.Origin, l.Color) }

func (l light) marshalConstructor() string {
	return fmt.Sprintf("newLight(%s, %s)", l.Origin.marshalConstructor(), l.Color.marshalConstructor())
}

func getLight(in mat4) (center vec3, color vec4) {
	return in[0].vec3, in[2]
}

type camera struct {
	Origin vec3 `json:"origin"`
	LookAt vec3 `json:"lookAt"`
}

func getCameraComponents(in mat4) (forward, right, up vec3) {
	return in[0].vec3, in[1].vec3, in[2].vec3
}
