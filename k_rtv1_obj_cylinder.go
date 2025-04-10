package main

func newCylinder(center1, center2 vec3, radius float, materialIndex int) mat4 {
	return newMat4(
		newVec4(CylinderType, float(materialIndex), radius, 0),
		newVec4(center1.x, center1.y, center1.z, 0),
		newVec4(center2.x, center2.y, center2.z, 0),
		newVec4(0, 0, 0, 0),
	)
}
