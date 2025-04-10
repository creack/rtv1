package main

// This file contains the recursive functions used for the raytracer.
// It compiles to both Go and Kage shader (after pre-processing).

func getThingDiffuse(thing mat4, recPoint vec3, materials MaterialsT) vec4 {
	if t := getThingType(thing); t == SphereType {
		return diffuseSphere(thing, recPoint, materials)
	} else if t == PlaneType {
		return diffusePlane(thing, recPoint, materials)
	} else if t == ConeType {
		return diffuseCone(thing, recPoint, materials)
	} else if t == CylinderType {
		return newVec4(0, 1, 0, 1) // TODO: Implement diffuseCylinder
		// return diffuseCylinder(thing, recPoint, materials)
	}

	return newVec4(1, 0, 1, 1) // Error color.
}

//rec:func:trace
func trace(cameraOrigin vec3, rayDir vec3, lights LightsT, things ThingsT, materials MaterialsT, ambientLight mat4, depth int, x, y int) vec4 {
	result := newVec4(0.1, 0.1, 0.1, 1) // Background color.
	closestThing, dist := intersection(cameraOrigin, rayDir, things, 0.001, -1)

	if dist == 0 {
		return result
	}

	var hitNormal vec3

	hitPoint := add3(cameraOrigin, scale3(rayDir, dist))
	if t := getThingType(closestThing); t == SphereType {
		result = diffuseSphere(closestThing, hitPoint, materials)
		center, radius, _ := getSphere(closestThing)
		hitNormal = scale3(sub3(hitPoint, center), 1/radius)
	} else if t == PlaneType {
		result = diffusePlane(closestThing, hitPoint, materials)
		hitNormal = normalPlane(closestThing, hitPoint)
	} else if t == ConeType {
		result = diffuseCone(closestThing, hitPoint, materials)
		hitNormal = normalCone(closestThing, hitPoint)
	} else if t == CylinderType {
		return newVec4(0, 1, 0, 1) // TODO.
	} else {
		return newVec4(1, 1, 0, 1) // Error color.
	}

	_, matAmbient, matDiffuse, matSpecular, matSpecularPower, matReflectiveIndex := getMaterial(materials, getThingMaterialIdx(closestThing))

	// Initialize the result with the ambient light.
	_, ambientLightColor, ambientLightIntensity := getLight(ambientLight)
	result = mul4(scale4(result, matAmbient), scale4(ambientLightColor, ambientLightIntensity))

	for i := 0; i < len(lights); i++ {
		light := lights[i]

		// Get the light fields from the object.
		lightOrigin, lightColor, lightIntensity := getLight(light)

		// Calculate the light direction and distance.
		lightDir := sub3(lightOrigin, hitPoint)
		lightDistance := length3(lightDir)
		lightDir = normalize3(lightDir)

		// Re-cast from the hit point to the light source.
		_, dist := intersection(hitPoint, lightDir, things, 0.001, lightDistance)
		if dist != 0 { // If we hit something, we don't see the light, so move forward.
			continue
		}

		// If we didn't hit anything, it means we have the light source in sight.

		// Diffuse lighting.
		diffFactor := max(0, dot3(hitNormal, lightDir))
		diffuse := scale4(getThingDiffuse(closestThing, hitPoint, materials), matDiffuse*diffFactor)

		// Specular lighting.
		viewDir := normalize3(scale3(rayDir, -1))
		reflectDir := reflect3(scale3(lightDir, -1), hitNormal)
		specFactor := pow(max(0, dot3(viewDir, reflectDir)), matSpecularPower)
		specular := scale4(lightColor, matSpecular*specFactor)

		// Combine diffuse and specular components.
		combined := add4(diffuse, specular)

		// Apply the light color and intensity.
		combined = scale4(mul4(combined, lightColor), lightIntensity)

		// Apply distance attenuation (inverse square law).
		attenuation := 1.0 / (lightDistance * lightDistance)
		combined = scale4(combined, attenuation)

		result = add4(result, combined)
	}

	_ = matReflectiveIndex
	//rec:if:depth
	if matReflectiveIndex > 0 && depth > 0 {
		reflectDir := reflect3(rayDir, hitNormal)
		//rec:rec-call:trace
		reflectColor := trace(hitPoint, reflectDir, lights, things, materials, ambientLight, depth-1, x, y)
		result = add4(result, scale4(reflectColor, matReflectiveIndex))
	}
	//rec:endif:depth

	return result
}

//rec:endfunc:trace
