package main

//kage:unit pixels

// This file is the main RTv1 logic. It compiles to both Go and Kage shader (after pre-processing).

// getThingType returns the type of the thing.
// By convention, it is stored in the x component of the first column of the mat4.
func getThingType(thing mat4) float {
	return thing[0].x
}

// getThingMaterial returns the material of the thing.
// By convention, it is stored in the y component of the first column of the mat4.
func getThingMaterialIdx(thing mat4) int {
	return int(thing[0].y)
}

// getMaterialColor returns the color of the material.
func getMaterialColor(materials MaterialsT, idx int) vec4 {
	color, _, _, _, _, _ := getMaterial(materials, idx)
	return color
}

const maxDepth = 15

const (
	SphereType   = 1
	PlaneType    = 2
	ConeType     = 3
	CylinderType = 4
)

// Fragment is the shader's entry point.
//
//nolint:revive // Unexported return is required by the shader API.
func Fragment(position vec4, _ vec2, _ vec4) vec4 {
	// "Localize" the uniform globals.
	width, height := int(Resolution.x), int(Resolution.y)

	cameraOrigin, cameraLookAt := UniCameraOrigin, UniCameraLookAt
	// cameraOrigin = newVec3(5*cos(0.5*Time), 0, 5*sin(0.5*Time))

	x := int(position.x)
	y := int(position.y)

	// Inject the scene constructors for the shader mode.
	// In Go mode, we use the global variables.
	//scene:things
	//scene:lights
	//scene:materials
	//scene:ambientLight

	cameraComponents := newCameraComponents(cameraOrigin, cameraLookAt)

	rayDir := initRay(width, height, x, y, cameraComponents)

	out := trace(cameraOrigin, rayDir, sceneLights, sceneObjects, sceneMaterials, ambientLight, maxDepth, x, y)

	return out
}
