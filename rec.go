package main

//rec:func:trace
func trace(camera mat4, rayDir vec3, lights LightsT, things ThingsT, depth int) vec4 {
	rayStart, _, _, _ := getCamera(camera) //nolint:dogsled // Expected.

	closestThing, dist := intersection(rayStart, rayDir, things)
	if dist == 0 {
		return newVec4(0, 0, 0, 1)
	}
	//rec:call:shade
	return shade(rayStart, rayDir, closestThing, dist, lights, things, depth)
}

//rec:endfunc:trace

//rec:func:shade
func shade(rayStart, rayDir vec3, closestThing mat4, dist float, lights LightsT, things ThingsT, depth int) vec4 {
	d := rayDir
	pos := add3(rayStart, scale3(d, dist))
	var normal vec3
	if getThingType(closestThing) == 1 {
		center, _, _, _ := getSphere(closestThing)
		normal = normalSphere(pos, center)
	} else {
		center, _, _, _ := getSphere(closestThing)
		normal = normalPlane(pos, center)
	}

	nd := dot3(normal, d)
	snd := scale3(normal, nd)
	s2nd := scale3(snd, 2)
	reflectDir := sub3(d, s2nd)

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
		var camera mat4
		camera[0] = newVec4(pos.x, pos.y, pos.z, 0)

		var reflected float
		if getThingType(thing) == 1 {
			reflected = reflectSphere(thing, pos)
		} else {
			reflected = reflectPlane(thing, pos)
		}
		return scale4(
			//rec:rec-call:trace
			trace(camera, rd, lights, things, depth+1),
			reflected,
		)
	}
	//rec:endif:depth
	return reflectedColor
}

//rec:endfunc:getReflectionColor
