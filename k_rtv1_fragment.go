package main

// Fragment is the shader's entry point.
//
//nolint:revive // Unexported return is required by the shader API.
func Fragment(position vec4, _ vec2, _ vec4) vec4 {
	// "Localize" the uniform globals.
	width, height := UniScreenWidth, UniScreenHeight
	cameraOrigin, cameraLookAt := UniCameraOrigin, UniCameraLookAt

	x := int(position.x)
	y := int(position.y)

	// Inject the scene constructors for the shader mode.
	// In Go mode, we use the global variables.
	//scene:things
	//scene:lights

	cameraComponents := newCameraComponents(cameraOrigin, cameraLookAt)

	rayDir := initRay(width, height, x, y, cameraComponents)

	out := trace(cameraOrigin, rayDir, sceneLights, sceneObjects, 0)

	return out
}
