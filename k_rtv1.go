package main

//kage:unit pixels

// This file is the main RTv1 logic. It compiles to both Go and Kage shader (after pre-processing).

// getThingType returns the type of the thing.
// By convention, it is stored in the z component of the second column of the mat4.
func getThingType(thing mat4) float {
	return thing[1].z
}

const maxDepth = 5

const (
	SphereType = 1
	PlaneType  = 2
	LightType  = 3
)
