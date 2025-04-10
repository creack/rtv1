package main

// p[1].xyz = center
// p[1].w = intensity
// p[1].xyzw = color
func newLight(center vec3, color vec4, intensity float) mat4 {
	return newMat4(
		newVec4(0, 0, 0, 0),
		newVec4(center.x, center.y, center.z, intensity),
		color,
		newVec4(0, 0, 0, 0),
	)
}

func getLight(in mat4) (center vec3, color vec4, intensity float) {
	return in[1].xyz, in[2], in[1].w
}
