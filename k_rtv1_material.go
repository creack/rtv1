package main

// p[0].x = type
// p[1].xyzw = color
// p[0].y = ambient
// p[0].z = diffuse
// p[0].w = specular
// p[2].x = specularPower
// p[2].y = reflectiveIndex
func newMaterial(mType int, color vec4, ambient, diffuse, specular, specularPower, reflectiveIndex float) mat4 {
	return newMat4(
		newVec4(float(mType), ambient, diffuse, specular),
		newVec4(color.x, color.y, color.z, color.w),
		newVec4(specularPower, reflectiveIndex, 0, 0),
		newVec4(0, 0, 0, 0),
	)
}
