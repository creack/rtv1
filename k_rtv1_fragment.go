package main

// Fragment is the shader's entry point.
//
//nolint:revive // Unexported return is required by the shader API.
func Fragment(position vec4, _ vec2, _ vec4) vec4 {
	width := UniScreenWidth
	height := UniScreenHeight
	x := int(position.x)
	y := int(position.y)

	things := ThingsT{
		newPlane(newVec3(0, 1.0, 0), 0, newVec4(1, 1, 1, 1)),
		newSphere(newVec3(0, 1, -0.25), 1, newVec4(1, 1, 0, 1)),
		newSphere(newVec3(-1.0, 0.5, 1.5), 0.5, newVec4(1, 0, 0, 1)),
	}

	lights := LightsT{
		// newLight(newVec3(-2.0, 2.5, 0), newVec4(1, 1, 1, 1)),
		newLight(newVec3(-2.0, 2.5, 0), newVec4(0.49, 0.07, 0.07, 1)),
		newLight(newVec3(1.5, 2.5, 1.5), newVec4(0.07, 0.07, 0.49, 1)),
		newLight(newVec3(1.5, 2.5, -1.5), newVec4(0.07, 0.49, 0.07, 1)),
		newLight(newVec3(0, 3.5, 0), newVec4(0.21, 0.21, 0.35, 1)),
	}
	//
	camera := newCamera(UniCameraOrigin, UniCameraLookAt)
	//
	rayDir := initRay(width, height, x, y, camera)
	out := trace(camera, rayDir, lights, things, 0)

	return out
}
