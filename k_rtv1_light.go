package main

// p[0].xyz = center
// p[0].w = intensity
// p[1].xyzw = color
func newLight(center vec3, color vec4, intensity float) mat4 {
	return newMat4(
		newVec4(center.x, center.y, center.z, intensity),
		newVec4(0, 0, LightType, 0),
		color,
		newVec4(0, 0, 0, 0),
	)
}

func getLight(in mat4) (center vec3, color vec4, intensity float) {
	return in[0].xyz, in[2], in[0].w
}

//var first int

// func addLight(thing mat4, pos, norm, rd vec3, col vec4, light mat4, things ThingsT) vec4 {
// 	lightPos, lightColor, _ := getLight(light)

// 	ldis := sub3(lightPos, pos)
// 	livec := normalize3(ldis)

// 	rayStart := pos
// 	rayDir := livec
// 	neatIsect := testRay(rayStart, rayDir, things)

// 	isInShadow := neatIsect != -1 && neatIsect <= length3(ldis)
// 	if isInShadow {
// 		return col
// 	}

// 	illum := dot3(livec, norm)
// 	lcolor := newVec4(0, 0, 0, 1) // defaultColor.
// 	if illum > 0 {
// 		lcolor = scale4(lightColor, illum)
// 	}

// 	specular := dot3(livec, normalize3(rd))
// 	scolor := newVec4(0, 0, 0, 1) // defaultColor.
// 	if specular > 0 {
// 		roughness := 0.
// 		if t := getThingType(thing); t == SphereType {
// 			roughness = roughnessSphere(thing, pos)
// 		} else if t == PlaneType {
// 			roughness = roughnessPlane(thing, pos)
// 		} else {
// 			roughness = -1.
// 		}
// 		scolor = scale4(lightColor, pow(specular, roughness))
// 	}
// 	var surfaceSpecular, surfaceDiffuse vec4
// 	if t := getThingType(thing); t == SphereType {
// 		surfaceSpecular = specularSphere(thing, pos)
// 		surfaceDiffuse = diffuseSphere(thing, pos)
// 	} else if t == PlaneType {
// 		surfaceSpecular = specularPlane(thing, pos)
// 		surfaceDiffuse = diffusePlane(thing, pos)
// 	}
// 	return add4(
// 		col,
// 		add4(
// 			mul4(surfaceDiffuse, lcolor),
// 			mul4(surfaceSpecular, scolor),
// 		),
// 	)
// }

// func getNaturalColor(thing mat4, pos, norm, rd vec3, lights LightsT, things ThingsT) vec4 {
// 	defaultColor0 := newVec4(0.1, 0.1, 0.1, 1)

// 	out := defaultColor0
// 	for i := 0; i < len(lights); i++ {
// 		out = addLight(thing, pos, norm, rd, out, lights[i], things)
// 	}
// 	return out
// }
