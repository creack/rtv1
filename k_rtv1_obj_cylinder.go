package main

// c[0].x = type - cylinder
// c[0].y = materialIdx
// c[0].z = radius
// c[1].xyz = center1
// c[2].xyz = center2
func newCylinder(center1, center2 vec3, radius float, materialIndex int) mat4 {
	return newMat4(
		newVec4(CylinderType, float(materialIndex), radius, 0),
		newVec4(center1.x, center1.y, center1.z, 0),
		newVec4(center2.x, center2.y, center2.z, 0),
		newVec4(0, 0, 0, 0),
	)
}

func getCylinder(in mat4) (center1, center2 vec3, radius float) {
	return in[1].xyz, in[2].xyz, in[0].z
}

func diffuseCylinder(thing mat4, pos vec3, materials MaterialsT) vec4 {
	_ = pos
	return getMaterialColor(materials, getThingMaterialIdx(thing))
}

func normalCylinder(thing mat4, pos vec3) vec3 {
	center1, center2, _ := getCylinder(thing)

	axis := sub3(center2, center1)
	axisDir := normalize3(axis)

	hitPointOnAxis := dot3(sub3(pos, center1), axisDir)

	pointOnAxis := add3(center1, scale3(axisDir, hitPointOnAxis))

	normal := sub3(pos, pointOnAxis)
	normal = normalize3(normal)

	return normal
}

func hitCylinder(rayStart, rayDir vec3, thing mat4, minDist, maxDist float) float {
	center1, center2, radius := getCylinder(thing)

	// Get axis information.
	axis := sub3(center2, center1)
	axisLength := length3(axis)
	axisDir := normalize3(axis)

	// Vector from ray start to center1.
	oc := sub3(rayStart, center1)

	// Project ray direction onto axis.
	rayDirDotAxis := dot3(rayDir, axisDir)

	// Perpendicular component of ray direction to axis.
	rayDirPerp := sub3(rayDir, scale3(axisDir, rayDirDotAxis))
	rayDirPerpLenSq := dot3(rayDirPerp, rayDirPerp)

	// Project oc onto axis.
	ocDotAxis := dot3(oc, axisDir)

	// Perpendicular component of oc to axis.
	ocPerp := sub3(oc, scale3(axisDir, ocDotAxis))

	// Set up quadratic equation coefficients.
	a := rayDirPerpLenSq
	if abs(a) < 1e-8 { // If ray is parallel to cylinder axis, no intersection.
		return 0
	}
	b := 2.0 * dot3(rayDirPerp, ocPerp)
	c := dot3(ocPerp, ocPerp) - radius*radius
	discriminant := b*b - 4.0*a*c

	if discriminant < 0.0 {
		return 0
	}

	// Find closest intersection.
	sqrtd := sqrt(discriminant)
	t1 := (-b - sqrtd) / (2.0 * a)
	t2 := (-b + sqrtd) / (2.0 * a)

	// Ensure t1 <= t2.
	if t1 > t2 {
		t1, t2 = t2, t1
	}

	t := t1
	// Check if first intersection is within bounds.
	hitPoint := add3(rayStart, scale3(rayDir, t))
	hitPointOnAxis := dot3(sub3(hitPoint, center1), axisDir)

	// If first hit is outside cylinder height, try second intersection.
	if hitPointOnAxis < 0 || hitPointOnAxis > axisLength || t < minDist || (maxDist != -1 && t > maxDist) {
		t = t2
		if t < minDist || (maxDist != -1 && t > maxDist) {
			return 0
		}

		hitPoint = add3(rayStart, scale3(rayDir, t))
		hitPointOnAxis = dot3(sub3(hitPoint, center1), axisDir)

		if hitPointOnAxis < 0 || hitPointOnAxis > axisLength {
			return 0
		}
	}

	return t

}
