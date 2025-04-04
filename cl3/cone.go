package main

import "math"

// Cone represents a cone in 3D space
type Cone struct {
	Apex     Vector3
	Base     Vector3
	Radius   float64
	Material Material
}

// Hit checks if a ray hits the cone
func (cone Cone) Hit(ray Ray, tMin, tMax float64) (bool, HitRecord) {
	// Vector from apex to base center
	axis := cone.Base.Sub(cone.Apex)
	axisLength := axis.Length()
	axisDir := axis.Normalize()

	// Calculate cone parameters
	cosTheta := axisLength / math.Sqrt(axisLength*axisLength+cone.Radius*cone.Radius)
	cosTheta2 := cosTheta * cosTheta

	// Vector from apex to ray origin
	oc := ray.Origin.Sub(cone.Apex)

	// Ray direction dot axis
	rdDotAxis := ray.Direction.Dot(axisDir)

	// Ray origin dot axis
	ocDotAxis := oc.Dot(axisDir)

	// Quadratic equation coefficients
	a := rdDotAxis*rdDotAxis - cosTheta2*ray.Direction.Dot(ray.Direction)
	b := 2.0 * (rdDotAxis*ocDotAxis - cosTheta2*ray.Direction.Dot(oc))
	c := ocDotAxis*ocDotAxis - cosTheta2*oc.Dot(oc)

	// Check discriminant
	discriminant := b*b - 4.0*a*c
	if discriminant < 0.0 {
		return false, HitRecord{}
	}

	// Calculate intersection
	sqrtd := math.Sqrt(discriminant)

	// Calculate both intersection points
	t0 := (-b - sqrtd) / (2.0 * a)
	t1 := (-b + sqrtd) / (2.0 * a)

	// Ensure t0 <= t1
	if t0 > t1 {
		t0, t1 = t1, t0
	}

	// Check both intersection points
	t := t0
	if t < tMin || t > tMax {
		t = t1
		if t < tMin || t > tMax {
			return false, HitRecord{}
		}
	}

	// Calculate intersection point
	hitPoint := ray.PointAt(t)

	// Check if intersection is within cone height (between apex and base)
	v := hitPoint.Sub(cone.Apex)
	projLen := v.Dot(axisDir)

	if projLen < 0.0 || projLen > axisLength {
		// Try other intersection
		t = t1
		if t < tMin || t > tMax {
			return false, HitRecord{}
		}

		hitPoint = ray.PointAt(t)
		v = hitPoint.Sub(cone.Apex)
		projLen = v.Dot(axisDir)

		if projLen < 0.0 || projLen > axisLength {
			return false, HitRecord{}
		}
	}

	// Calculate normal
	// Project hit point onto axis
	axisPoint := cone.Apex.Add(axisDir.Mul(projLen))

	// Get perpendicular component
	perpComp := hitPoint.Sub(axisPoint)

	// For cone normal, we need to account for slant
	normal := perpComp.Add(axisDir.Mul(-cone.Radius / axisLength))
	normal = normal.Normalize()

	return true, HitRecord{
		T:        t,
		Point:    hitPoint,
		Normal:   normal,
		Material: cone.Material,
	}
}
