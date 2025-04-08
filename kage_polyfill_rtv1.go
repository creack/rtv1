package main

// This file is a wrapper for Kage. It mirror the shaderlib_rtv1.kage file to allow
// for the shader to be compile in Go.
// In this part, we have the RTv1 specific functions and types.

type ThingsT [3]mat4

type LightsT [4]mat4

//nolint:unparam // Keeping for reference.
func getSphere(in mat4) (center vec3, radius, radius2 float, color vec4) {
	return in[0].vec3, in[0].w, in[1].w, in[2]
}

func getPlane(in mat4) (center vec3, offset, roughness float, color vec4) {
	return in[0].vec3, in[0].w, in[1].x, in[2]
}

func getLight(in mat4) (center vec3, color vec4) {
	return in[0].vec3, in[2]
}

func getCameraComponents(in mat4) (forward, right, up vec3) {
	return in[0].vec3, in[1].vec3, in[2].vec3
}
