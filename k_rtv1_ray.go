package main

func intersection(rayStart, rayDir vec3, things ThingsT) (closestThing mat4, closest float) {
	closest = -1.

	for i := 0; i < len(things); i++ {
		if dist := intersect(rayStart, rayDir, things[i]); dist > 0 {
			if closest == -1 || dist < closest {
				closestThing = things[i]
				closest = dist
			}
		}
	}
	if closest == -1. {
		closest = 0.
	}

	return closestThing, closest
}

func testRay(rayStart, rayDir vec3, things ThingsT) float {
	_, dist := intersection(rayStart, rayDir, things)
	if dist != 0 {
		return dist
	}
	return -1
}

func initRay(width, height, x, y int, cameraComponents mat4) vec3 {
	recenterX := (float(x) - float(width)/2.0) / 2.0 / float(width)
	recenterY := -(float(y) - float(height)/2.0) / 2.0 / float(height)

	forward, right, up := getCameraComponents(cameraComponents)

	return normalize3(add3(forward, add3(scale3(right, recenterX), scale3(up, recenterY))))
}
