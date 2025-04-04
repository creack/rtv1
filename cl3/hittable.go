package main

// HitRecord contains information about a ray hit
type HitRecord struct {
	T        float64
	Point    Vector3
	Normal   Vector3
	Material Material
}

// Hittable is an interface for objects that can be hit by a ray
type Hittable interface {
	Hit(ray Ray, tMin, tMax float64) (bool, HitRecord)
}

// HittableList is a collection of hittable objects
type HittableList struct {
	Objects []Hittable
}

// Hit checks if a ray hits any object in the list
func (h HittableList) Hit(ray Ray, tMin, tMax float64) (bool, HitRecord) {
	hitAnything := false
	closestSoFar := tMax
	var rec HitRecord

	for _, object := range h.Objects {
		hit, tempRec := object.Hit(ray, tMin, closestSoFar)
		if hit {
			hitAnything = true
			closestSoFar = tempRec.T
			rec = tempRec
		}
	}

	return hitAnything, rec
}
