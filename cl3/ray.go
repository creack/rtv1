package main

// Ray represents a ray with origin and direction
type Ray struct {
	Origin    Vector3
	Direction Vector3
}

// PointAt returns a point along the ray at the given parameter t
func (r Ray) PointAt(t float64) Vector3 {
	return r.Origin.Add(r.Direction.Mul(t))
}
