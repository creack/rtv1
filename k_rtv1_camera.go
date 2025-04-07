package main

func newCamera(camStart, camLookAt vec3) mat4 {
	forward := UniCameraForward
	right := UniCameraRight
	up := UniCameraUp
	return newMat4(
		newVec4(camStart.x, camStart.y, camStart.z, 0),
		newVec4(forward.x, forward.y, forward.z, 0),
		newVec4(right.x, right.y, right.z, 0),
		newVec4(up.x, up.y, up.z, 0),
	)
}
