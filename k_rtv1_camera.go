package main

// newCameraComponents creates a camera transformation matrix.
// It returns a mat4 that contains the camera's forward, right, and up vectors.
//
// forward: p[0].xyz
// right: p[1].xyz
// up: p[2].xyz
func newCameraComponents(camStart, camLookAt vec3) mat4 {
	worldUp := newVec3(0, 1, 0)

	forward := normalize3(sub3(camLookAt, camStart))
	right := normalize3(cross3(worldUp, forward))
	up := normalize3(cross3(forward, right))

	return newMat4(
		newVec4(forward.x, forward.y, forward.z, 0),
		newVec4(right.x, right.y, right.z, 0),
		newVec4(up.x, up.y, up.z, 0),
		newVec4(0, 0, 0, 0),
	)
}

// calculatePitch is a helper function to calculate the pitch angle
// of the camera.
// Only used in the ebiten update move logic to restrict
// up/ddown rotation of the camera.
func calculatePitch(origin, lookAt vec3) float {
	// Get direction vector from origin to lookAt
	direction := sub3(lookAt, origin)

	// Get horizontal distance (length of direction projected on XZ plane)
	horizontalDist := sqrt(direction.x*direction.x + direction.z*direction.z)

	// Calculate pitch angle (vertical angle)
	// This will return negative for looking down, positive for looking up
	pitch := atan2(direction.y, horizontalDist)

	return pitch
}
