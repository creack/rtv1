package main

import (
	"math"
)

// Sphere represents a sphere in 3D space
type Sphere struct {
	Center   Vector3
	Radius   float64
	Material Material
}

var first int

// Hit checks if a ray hits the sphere
func (s Sphere) Hit(ray Ray, tMin, tMax float64) (bool, HitRecord) {
	oc := ray.Origin.Sub(s.Center)
	a := ray.Direction.Dot(ray.Direction)
	b := 2.0 * oc.Dot(ray.Direction)
	c := oc.Dot(oc) - s.Radius*s.Radius
	discriminant := b*b - 4*a*c

	// if !first {
	// 	first = true
	// 	fmt.Printf("%.3g,%.3g,%.3g | %.3g,%.3g,%.3g | %.3g,%.3g,%.3g | %.3g,%.3g,%.3g | %.3g | %.3g | %.3g | %.3g\n",
	// 		ray.Origin.X, ray.Origin.Y, ray.Origin.Z,
	// 		ray.Direction.X, ray.Direction.Y, ray.Direction.Z,
	// 		s.Center.X, s.Center.Y, s.Center.Z,
	// 		oc.X, oc.Y, oc.Z,
	// 		a, b, c, discriminant)
	// }

	if discriminant < 0 {
		return false, HitRecord{}
	}

	// Find the nearest root that lies in the acceptable range
	sqrtd := math.Sqrt(discriminant)
	root := (-b - sqrtd) / (2 * a)
	if root < tMin || root > tMax {
		root = (-b + sqrtd) / (2 * a)
		if root < tMin || root > tMax {
			return false, HitRecord{}
		}
	}

	t := root
	// point := ray.Origin.Add(ray.Direction.Mul(t))
	point := ray.PointAt(t)
	normal := point.Sub(s.Center).Div(s.Radius)

	return true, HitRecord{
		T:        t,
		Point:    point,
		Normal:   normal,
		Material: s.Material,
	}
}
