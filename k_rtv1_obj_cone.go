package main

// c[0].x = type - cone
// c[0].y = material index
// c[0].z = radius
// c[0].w = radius^2
// c[1].xyz = base center
// c[2].xyz = apex center
func newCone(base, apex vec3, radius float, materialIdx int) mat4 {
	return newMat4(
		newVec4(ConeType, float(materialIdx), radius, radius*radius),
		newVec4(base.x, base.y, base.z, 0),
		newVec4(apex.x, apex.y, apex.z, 0),
		newVec4(0, 0, 0, 0),
	)
}

func getCone(in mat4) (base, apex vec3, radius, radius2 float) {
	return in[1].xyz, in[2].xyz, in[0].z, in[0].w
}
func diffuseCone(thing mat4, pos vec3, materials MaterialsT) vec4 {
	_ = pos
	return getMaterialColor(materials, getThingMaterialIdx(thing))
}

func normalCone(thing mat4, pos vec3) vec3 {
	base, apex, radius, _ := getCone(thing)

	axis := sub3(base, apex)
	axisLength := length3(axis)
	axisDir := normalize3(axis)

	v := sub3(pos, apex)
	projLen := dot3(v, axisDir)

	// Project hit point onto axis.
	axisPoint := add3(apex, scale3(axisDir, projLen))

	perpComp := sub3(pos, axisPoint)

	normal := add3(perpComp, scale3(axisDir, -radius/axisLength))
	normal = normalize3(normal)

	return normal
}

func hitCone(rayStart, rayDir vec3, thing mat4, minDist, maxDist float) float {
	base, apex, _, radius2 := getCone(thing)

	// Vector from apex to base center.
	axis := sub3(base, apex)
	axisLength := length3(axis)
	axisDir := normalize3(axis)

	// Calculate the cone parameters.
	cosTheta := axisLength / sqrt(radius2+axisLength*axisLength)
	cosTheta2 := cosTheta * cosTheta

	// Vector from apex to ray origin.
	oc := sub3(rayStart, apex)

	// Ray direction dot axis.
	rdDotAxis := dot3(rayDir, axisDir)

	// Ray origin dot axis.
	ocDotAxis := dot3(oc, axisDir)

	// Quadratic equation coefficients.
	a := rdDotAxis*rdDotAxis - cosTheta2*dot3(rayDir, rayDir)
	b := 2.0 * (rdDotAxis*ocDotAxis - cosTheta2*dot3(rayDir, oc))
	c := ocDotAxis*ocDotAxis - cosTheta2*dot3(oc, oc)

	// Check discriminant.
	discriminant := b*b - 4.0*a*c
	if discriminant < 0.0 {
		return 0
	}

	// Calculate intersection.
	sqrtd := sqrt(discriminant)

	// Calculate both intersection points.
	t0 := (-b - sqrtd) / (2.0 * a)
	t1 := (-b + sqrtd) / (2.0 * a)

	// Ensure t0 <= t1.
	if t0 > t1 {
		t0, t1 = t1, t0
	}

	// Check both intersection points.
	t := t0
	if t < minDist || (maxDist != -1 && t > maxDist) {
		t = t1
		if t < minDist || (maxDist != -1 && t > maxDist) {
			return 0
		}
	}

	// Calculate intersection point.
	hitPoint := add3(rayStart, scale3(rayDir, t))

	// Check if intersection is within cone height (between apex and base).
	v := sub3(hitPoint, apex)
	projLen := dot3(v, axisDir)

	if projLen < 0.0 || projLen > axisLength {
		// Try other intersection.
		t = t1
		if t < minDist || (maxDist != -1 && t > maxDist) {
			return 0
		}

		hitPoint = add3(rayStart, scale3(rayDir, t))
		v = sub3(hitPoint, apex)
		projLen = dot3(v, axisDir)

		if projLen < 0.0 || projLen > axisLength {
			return 0
		}
	}

	return t
}
