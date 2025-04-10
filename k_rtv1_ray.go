package main

func intersect(rayStart, rayDir vec3, thing mat4, minDist, maxDist float) float {
	if t := getThingType(thing); t == SphereType {
		return hitSphere(rayStart, rayDir, thing, minDist, maxDist)
	} else if t == PlaneType {
		return hitPlane(rayStart, rayDir, thing, minDist, maxDist)
	} else if t == ConeType {
		return hitCone(rayStart, rayDir, thing, minDist, maxDist)
	} else if t == CylinderType {
		return -1
		// return hitCylinder(rayStart, rayDir, thing, minDist, maxDist)
	}
	return -1.
}

func intersection(rayStart, rayDir vec3, things ThingsT, minDist, maxDist float) (closestThing mat4, closest float) {
	closest = maxDist
	hitSomething := false
	for i := 0; i < len(things); i++ {
		dist := intersect(rayStart, rayDir, things[i], minDist, closest)
		hit := dist != 0
		if hit {
			hitSomething = true
			closest = dist
			closestThing = things[i]
		}

	}
	if !hitSomething {
		closest = 0.
	}

	return closestThing, closest
}

func initRay(width, height, x, y int, cameraComponents mat4) vec3 {
	// Hard-coded FOV for now.
	const FOV = 45.0

	// Calculcate the viewplane.
	aspectRatio := float(width) / float(height)
	theta := FOV * pi / 180.0
	halfHeight := tan(theta / 2.0)
	halfWidth := aspectRatio * halfHeight

	// Get the camera vectors.
	forward, right, up := getCameraComponents(cameraComponents)

	u := float(x) / float(width)
	v := 1.0 - float(y)/float(height)

	dir := scale3(right, u*2.0*halfWidth-halfWidth)
	dir = add3(dir, scale3(up, v*2.0*halfHeight-halfHeight))
	dir = sub3(dir, forward)
	dir = normalize3(dir)

	return dir // Return the calculated direction vector
}
