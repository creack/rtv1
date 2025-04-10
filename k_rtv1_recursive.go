package main

// This file contains the recursive functions used for the raytracer.
// It compiles to both Go and Kage shader (after pre-processing).

func getThingDiffuse(thing mat4, recPoint vec3, materials MaterialsT) vec4 {
	if t := getThingType(thing); t == SphereType {
		return diffuseSphere(thing, recPoint, materials)
	} else if t == PlaneType {
		return diffusePlane(thing, recPoint, materials)
	} else {
		return newVec4(1, 0, 1, 1) // Error color.
	}
}

func trace(cameraOrigin vec3, rayDir vec3, lights LightsT, things ThingsT, materials MaterialsT, depth int, x, y int) vec4 {
	result := newVec4(0.1, 0.1, 0.1, 1) // Background color.
	closestThing, dist := intersection(cameraOrigin, rayDir, things, 0.001, -1)

	if dist == 0 {
		return result
	}

	var recNormal vec3

	recPoint := add3(cameraOrigin, scale3(rayDir, dist))
	if t := getThingType(closestThing); t == SphereType {
		result = diffuseSphere(closestThing, recPoint, materials)
		center, radius, _, _ := getSphere(closestThing)
		recNormal = scale3(sub3(recPoint, center), 1/radius)
	} else if t == PlaneType {
		result = diffusePlane(closestThing, recPoint, materials)
		recNormal = normalPlane(closestThing, recPoint)
	} else {
		return newVec4(1, 1, 0, 1) // Error color.
	}

	const materialRedAmbient = 0.1
	getMaterial(materials, 0)
	ambientLight := newVec4(0.05, 0.05, 0.05, 1)
	result = mul4(scale4(result, materialRedAmbient), ambientLight)

	for i := 0; i < len(lights); i++ {
		light := lights[i]

		lightOrigin, lightColor, lightIntensity := getLight(light)

		lightDir := sub3(lightOrigin, recPoint)
		lightDistance := length3(lightDir)
		lightDir = normalize3(lightDir)

		// Re-cast from the hit point to the light source.
		_, dist := intersection(recPoint, lightDir, things, 0.001, lightDistance)
		if dist != 0 { // If we hit something, we don't see the light, so move forward.
			continue
		}

		// If we didn't hit anything, it means we have the light source in sight.

		const materialRedDiffuse = 0.7
		diffFactor := max(0, dot3(recNormal, lightDir))
		diffuse := scale4(getThingDiffuse(closestThing, recPoint, materials), materialRedDiffuse*diffFactor)

		viewDir := normalize3(scale3(rayDir, -1))
		reflectDir := reflect3(scale3(lightDir, -1), recNormal)
		const materalRedSpecular = 0.5
		const materalRedSpecularPower = 32
		specFactor := pow(max(0, dot3(viewDir, reflectDir)), materalRedSpecularPower)
		specular := scale4(lightColor, materalRedSpecular*specFactor)
		_ = specular

		combined := add4(diffuse, specular)
		combined = scale4(mul4(combined, lightColor), lightIntensity)

		attenuation := 1.0 / (lightDistance * lightDistance)
		combined = scale4(combined, attenuation)
		result = add4(result, combined)
	}

	return result
}
