package main

// type: p[0].x - plane
// materialIdx: p[0].y
// isCheckerBoard: p[0].z
// checkerSize: p[0].w
// center: p[1].xyz
// normal: p[2].xyz
func newPlane(center, normal vec3, isCheckerBoard bool, checkerSize float, materialIdx int) mat4 {
	isCheckerBoardFloat := 0.0
	if isCheckerBoard {
		isCheckerBoardFloat = 1.0
	}
	return newMat4(
		newVec4(PlaneType, float(materialIdx), isCheckerBoardFloat, checkerSize),
		newVec4(center.x, center.y, center.z, isCheckerBoardFloat),
		newVec4(normal.x, normal.y, normal.z, checkerSize),
		newVec4(0, 0, 0, 0),
	)
}

func diffusePlane(thing mat4, pos vec3, materials MaterialsT) vec4 {
	_, _, isCheckerboard, checkerSize := getPlane(thing) //nolint:dogsled // Expected.
	color := getMaterialColor(materials, getThingMaterialIdx(thing))
	if !isCheckerboard {
		return color
	}

	if (int(floor(pos.z/checkerSize))+int(floor(pos.x/checkerSize)))%2 != 0 {
		return newVec4(0.1, 0.1, 0.1, 1)
	}
	return color
}

func normalPlane(thing mat4, pos vec3) vec3 {
	_ = pos
	_, pNorm, _, _ := getPlane(thing) //nolint:dogsled // Expected.
	//log.Println(pNorm)
	return pNorm
}

func hitPlane(rayStart, rayDir vec3, thing mat4, minDist, maxDist float) float {
	pPos, pNorm, _, _ := getPlane(thing) //nolint:dogsled // Expected.

	denom := dot3(rayDir, pNorm)
	if abs(denom) < 1e-6 {
		return 0
	}

	dist := dot3(sub3(pPos, rayStart), pNorm) / denom
	if dist < minDist || (maxDist != -1 && dist > maxDist) {
		return 0
	}

	return dist
}
