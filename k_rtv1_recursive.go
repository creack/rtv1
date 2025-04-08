package main

// This file contains the recursive functions used for the raytracer.
// It compiles to both Go and Kage shader (after pre-processing).

//rec:func:trace
func trace(cameraOrigin vec3, rayDir vec3, lights LightsT, things ThingsT, depth int) vec4 {
	closestThing, dist := intersection(cameraOrigin, rayDir, things)
	if dist == 0 {
		return newVec4(0, 0, 0, 1)
	}
	//rec:call:shade
	return shade(cameraOrigin, rayDir, closestThing, dist, lights, things, depth)
}

//rec:endfunc:trace

//rec:func:shade
func shade(rayStart, rayDir vec3, closestThing mat4, dist float, lights LightsT, things ThingsT, depth int) vec4 {
	pos := add3(rayStart, scale3(rayDir, dist))
	var normal vec3
	if t := getThingType(closestThing); t == SphereType {
		center, _, _, _ := getSphere(closestThing)
		normal = normalSphere(pos, center)
	} else if t == PlaneType {
		center, _, _, _ := getPlane(closestThing)
		normal = normalPlane(pos, center)
	}

	nd := dot3(normal, rayDir)
	snd := scale3(normal, nd)
	s2nd := scale3(snd, 2)
	reflectDir := sub3(rayDir, s2nd)

	backgroundColor0 := newVec4(0, 0, 0, 0)

	naturalColor := add4(
		backgroundColor0,
		getNaturalColor(closestThing, pos, normal, reflectDir, lights, things),
	)

	//rec:call:getReflectionColor
	reflectedColor := getReflectionColor(
		closestThing,
		pos,
		reflectDir,
		lights,
		things,
		depth,
	)

	return add4(naturalColor, reflectedColor)
}

//rec:endfunc:shade

//rec:func:getReflectionColor
func getReflectionColor(thing mat4, pos, rd vec3, lights LightsT, things ThingsT, depth int) vec4 {
	grayColor := newVec4(0.3, 0.3, 0.3, 1)
	reflectedColor := grayColor

	// In shader mode, to avoid recursion, we exclude this block from the final depth.
	//rec:if:depth
	if depth+1 < maxDepth {
		var reflected float
		if getThingType(thing) == 1 {
			reflected = reflectSphere(thing, pos)
		} else {
			reflected = reflectPlane(thing, pos)
		}
		return scale4(
			//rec:rec-call:trace
			trace(pos, rd, lights, things, depth+1),
			reflected,
		)
	}
	//rec:endif:depth
	return reflectedColor
}

//rec:endfunc:getReflectionColor
