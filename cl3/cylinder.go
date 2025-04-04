package main

import "math"

// Cylinder represents a cylinder in 3D space
type Cylinder struct {
	Center1  Vector3
	Center2  Vector3
	Radius   float64
	Material Material
}

// Hit checks if a ray hits the cylinder
func (cyl Cylinder) Hit(ray Ray, tMin, tMax float64) (bool, HitRecord) {
	// Get axis information
	axis := cyl.Center2.Sub(cyl.Center1)
	axisLength := axis.Length()
	axisDir := axis.Normalize()

	// Vector from ray origin to cylinder base center
	oc := ray.Origin.Sub(cyl.Center1)

	// Project ray direction onto axis
	rayDirDotAxis := ray.Direction.Dot(axisDir)

	// Perpendicular component of ray direction to axis
	rayDirPerp := ray.Direction.Sub(axisDir.Mul(rayDirDotAxis))
	rayDirPerpLenSq := rayDirPerp.Dot(rayDirPerp)

	// Project oc onto axis
	ocDotAxis := oc.Dot(axisDir)

	// Perpendicular component of oc to axis
	ocPerp := oc.Sub(axisDir.Mul(ocDotAxis))

	// Set up quadratic equation coefficients
	a := rayDirPerpLenSq

	// If ray is parallel to cylinder axis, no intersection
	if math.Abs(a) < 1e-8 {
		return false, HitRecord{}
	}

	b := 2.0 * rayDirPerp.Dot(ocPerp)
	c := ocPerp.Dot(ocPerp) - cyl.Radius*cyl.Radius

	discriminant := b*b - 4.0*a*c

	if discriminant < 0.0 {
		return false, HitRecord{}
	}

	// Find closest intersection
	sqrtd := math.Sqrt(discriminant)
	t1 := (-b - sqrtd) / (2.0 * a)
	t2 := (-b + sqrtd) / (2.0 * a)

	// Ensure t1 <= t2
	if t1 > t2 {
		t1, t2 = t2, t1
	}

	t := t1
	// Check if first intersection is within bounds
	hitPoint := ray.PointAt(t)
	hitPointOnAxis := hitPoint.Sub(cyl.Center1).Dot(axisDir)

	// If first hit is outside cylinder height, try second intersection
	if hitPointOnAxis < 0 || hitPointOnAxis > axisLength || t < tMin || t > tMax {
		t = t2
		if t < tMin || t > tMax {
			return false, HitRecord{}
		}

		hitPoint = ray.PointAt(t)
		hitPointOnAxis = hitPoint.Sub(cyl.Center1).Dot(axisDir)

		if hitPointOnAxis < 0 || hitPointOnAxis > axisLength {
			return false, HitRecord{}
		}
	}

	// Find point on axis
	pointOnAxis := cyl.Center1.Add(axisDir.Mul(hitPointOnAxis))

	// Calculate normal vector (from axis to hit point)
	normal := hitPoint.Sub(pointOnAxis).Normalize()

	return true, HitRecord{
		T:        t,
		Point:    hitPoint,
		Normal:   normal,
		Material: cyl.Material,
	}
}
