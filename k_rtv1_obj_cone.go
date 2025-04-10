package main

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
	base, apex, _, _ := getCone(thing)
	// log.Println(base)
	// log.Println(apex)
	// log.Println(pos)

	d := sub3(apex, base)
	v := sub3(pos, base)
	n := cross3(d, v)
	n = normalize3(n)

	return n
}

func hitCone(rayStart, rayDir vec3, thing mat4, minDist, maxDist float) float {
	base, apex, _, radius2 := getCone(thing)
	// apex = apex + (rayDir * 0.0001)
	// log.Println(apex)
	// log.Println(rayStart)
	// log.Println(rayDir)
	// log.Println(base)
	// log.Println(radius2)

	d := sub3(apex, base)
	f := dot3(d, d) - radius2
	if f < 0 {
		return 0
	}

	m := dot3(sub3(rayStart, base), d) / f
	n := dot3(rayDir, d) / f

	a := dot3(rayDir, rayDir) - n*n
	b := dot3(rayDir, d) - m*n
	c := dot3(sub3(rayStart, base), rayDir) - m*f

	discriminant := b*b - a*c
	if discriminant < 0 {
		return 0
	}

	sqrtd := sqrt(discriminant)
	root := (-b - sqrtd) / a
	if root < minDist || (maxDist != -1 && root > maxDist) {
		root = (-b + sqrtd) / a
		if root < minDist || (maxDist != -1 && root > maxDist) {
			return 0
		}
	}

	return root
}
