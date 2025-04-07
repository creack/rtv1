package main

// center: s[0].xyz
// radius: s[0].w
// radius2: s[1].w
// roughness: s[1].x
// color: s[2]
// specular: s[3]
func newSphere(center vec3, radius float, col vec4) mat4 {
	return newMat4(
		newVec4(center.x, center.y, center.z, radius),
		newVec4(250., 0, SphereType, radius*radius),
		col,
		newVec4(0.3, 0.3, 0.3, 1),
	)
}

func reflectSphere(thing mat4, pos vec3) float {
	_ = pos
	_ = thing
	return 0.7
}

func diffuseSphere(thing mat4, pos vec3) vec4 {
	_ = pos
	_, _, _, col := getSphere(thing) //nolint:dogsled // Expected.
	return col
}

func specularSphere(thing mat4, pos vec3) vec4 {
	_ = pos
	return thing[3]
}

func roughnessSphere(thing mat4, pos vec3) float {
	_ = pos
	return thing[1].x
}

func normalSphere(rayStart, center vec3) vec3 {
	return normalize3(sub3(rayStart, center))
}

func hitSphere(rayStart, rayDir vec3, thing mat4) float {
	sphereCenter, _, sphereRadius2, _ := getSphere(thing)

	eo := sub3(sphereCenter, rayStart)
	v := dot3(eo, rayDir)

	disc := sphereRadius2 - (dot3(eo, eo) - v*v)
	if disc >= 0 {
		if dist := v - sqrt(disc); dist >= 0 {
			return dist
		}
	}
	return 0
}

func intersect(rayStart, rayDir vec3, thing mat4) float {
	if t := getThingType(thing); t == SphereType {
		return hitSphere(rayStart, rayDir, thing)
	} else if t == PlaneType {
		return hitPlane(rayStart, rayDir, thing)
	}
	return -1.
}
