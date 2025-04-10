package main

import "math"

// Camera represents a simple camera
type Camera struct {
	Origin      Vector3
	LookAt      Vector3
	Up          Vector3
	FOV         float64
	AspectRatio float64
}

// GetRay returns a ray from the camera to the given screen coordinates (u,v)
func (c Camera) GetRay(u, v float64) Ray {
	forward := c.Origin.Sub(c.LookAt).Normalize()
	right := c.Up.Cross(forward).Normalize()
	up := forward.Cross(right)

	// Calculate viewplane
	theta := c.FOV * math.Pi / 180.0
	half_height := math.Tan(theta / 2.0)
	half_width := c.AspectRatio * half_height

	// Calculate ray direction
	direction := right.Mul(u*2.0*half_width - half_width)
	direction = direction.Add(up.Mul(v*2.0*half_height - half_height))
	direction = direction.Sub(forward)
	direction = direction.Normalize()

	return Ray{c.Origin, direction}
}

// Default camera vectors
var (
	CameraPosition Vector3 = Vector3{0, 0, 5} // Set initial position farther back
	CameraForward  Vector3
	CameraRight    Vector3
	CameraUp       Vector3
)

var (
	Yaw   float64 = -math.Pi / 2 // Start looking in -Z direction
	Pitch float64 = 0
	Roll  float64 = 0
)

// UpdateCameraVectors updates the camera orientation vectors based on yaw and pitch
func UpdateCameraVectors() {
	// Calculate the new front vector
	CameraForward = Vector3{
		X: math.Cos(Yaw) * math.Cos(Pitch),
		Y: math.Sin(Pitch),
		Z: math.Sin(Yaw) * math.Cos(Pitch),
	}.Normalize()

	// Recalculate the right and up vectors
	worldUp := Vector3{0, 1, 0}
	CameraRight = worldUp.Cross(CameraForward).Normalize()
	CameraUp = CameraForward.Cross(CameraRight).Normalize()
}

// CreateCamera initializes a camera with default settings
func CreateCamera(screenWidth, screenHeight int) Camera {
	// Make sure camera vectors are initialized
	UpdateCameraVectors()

	// Create camera
	return Camera{
		Origin:      CameraPosition,
		LookAt:      CameraPosition.Add(CameraForward),
		Up:          CameraUp,
		FOV:         45.0,
		AspectRatio: float64(screenWidth) / float64(screenHeight),
	}
}
