package main

import "math"

// Plane represents an infinite plane in 3D space
type Plane struct {
	Point          Vector3
	Normal         Vector3
	Material       Material
	IsCheckerboard bool
	CheckerSize    float64
}

// Hit checks if a ray hits the plane
func (p Plane) Hit(ray Ray, tMin, tMax float64) (bool, HitRecord) {
	denom := ray.Direction.Dot(p.Normal)
	if math.Abs(denom) < 1e-6 {
		return false, HitRecord{} // Ray is parallel to plane
	}

	t := p.Point.Sub(ray.Origin).Dot(p.Normal) / denom
	if t < tMin || t > tMax {
		return false, HitRecord{}
	}

	hitPoint := ray.PointAt(t)
	material := p.Material

	// Apply checkerboard pattern if enabled
	if p.IsCheckerboard {
		// Calculate checkerboard pattern based on x and z coordinates
		x := math.Floor(hitPoint.X / p.CheckerSize)
		z := math.Floor(hitPoint.Z / p.CheckerSize)

		// If sum is even, use alternate color (black)
		if math.Mod(math.Abs(x+z), 2.0) < 0.5 {
			// Create a copy of the material with modified color
			material.Color = Color{0.1, 0.1, 0.1} // Dark color for alternate squares
		}
	}

	return true, HitRecord{
		T:        t,
		Point:    hitPoint,
		Normal:   p.Normal,
		Material: material,
	}
}
