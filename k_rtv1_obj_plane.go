package main

// center: p[0].xyz
// offset: p[0].w
// roughness: p[1].x
// color: p[2]
// specular: p[3]
func newPlane(center vec3, offset float, col vec4) mat4 {
	return newMat4(
		newVec4(center.x, center.y, center.z, offset),
		newVec4(150, 0, PlaneType, 0),
		col,
		newVec4(1, 1, 1, 1),
	)
}

func reflectPlane(thing mat4, pos vec3) float {
	_ = thing
	if (int(floor(pos.z))+int(floor(pos.x)))%2 != 0 {
		return 0.1
	}
	return 0.7
}

func diffusePlane(thing mat4, pos vec3) vec4 {
	if (int(floor(pos.z))+int(floor(pos.x)))%2 != 0 {
		return thing[2]
	}
	return newVec4(0, 0, 0, 0)
}

func specularPlane(thing mat4, pos vec3) vec4 {
	_ = pos
	return thing[3]
}

func roughnessPlane(thing mat4, pos vec3) float {
	_ = pos
	return thing[1].x
}

func normalPlane(rayStart, center vec3) vec3 {
	_ = rayStart
	return center
}

func hitPlane(rayStart, rayDir vec3, thing mat4) float {
	ppos, _, _, _ := getSphere(thing) //nolint:dogsled // Expected.
	denom := dot3(ppos, rayDir)
	if denom > 0 {
		return 0
	}
	dist := (dot3(ppos, rayStart) + thing[0].w) / (-denom)
	return dist
}
