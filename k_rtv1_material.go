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

func getMaterial(materials MaterialsT, idx int) (color vec4, ambient, diffuse, specular, specularPower, reflectiveIndex float) {
	m := materials[idx]
	color = m[1]
	ambient = m[0].y
	diffuse = m[0].z
	specular = m[0].w
	specularPower = m[2].x
	reflectiveIndex = m[2].y
	return color, ambient, diffuse, specular, specularPower, reflectiveIndex
}
