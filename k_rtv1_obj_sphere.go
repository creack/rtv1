package main

// type: s[0].x - sphere
// materialIdx: s[0].y
// radius: s[0].z
// radius^2: s[0].w
// center: s[1].xyz
func newSphere(center vec3, radius float, materialIdx int) mat4 {
	return newMat4(
		newVec4(SphereType, float(materialIdx), radius, radius*radius),
		newVec4(center.x, center.y, center.z, 0),
		newVec4(0, 0, 0, 0),
		newVec4(0, 0, 0, 0),
	)
}

func getSphere(in mat4) (center vec3, radius, radius2 float) {
	return newVec3(in[1].x, in[0].y+(cos(Time)), in[0].z), in[0].z, in[0].w
	//return in[1].xyz, in[0].z, in[0].w
}

func diffuseSphere(thing mat4, pos vec3, materials MaterialsT) vec4 {
	_ = pos
	return getMaterialColor(materials, getThingMaterialIdx(thing))
}

func hitSphere(rayStart, rayDir vec3, thing mat4, minDist, maxDist float) float {
	sphereCenter, _, sphereRadius2 := getSphere(thing)

	oc := sub3(rayStart, sphereCenter)
	a := dot3(rayDir, rayDir)
	b := 2.0 * dot3(oc, rayDir)
	c := dot3(oc, oc) - sphereRadius2
	discriminant := b*b - 4*a*c

	if discriminant < 0 {
		return 0
	}

	sqrtd := sqrt(discriminant)
	root := (-b - sqrtd) / (2 * a)
	if root < minDist || (maxDist != -1 && root > maxDist) {
		root = (-b + sqrtd) / (2 * a)
		if root < minDist || (maxDist != -1 && root > maxDist) {
			return 0
		}
	}

	return root
}
