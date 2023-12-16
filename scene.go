package main

import (
	"go.creack.net/rtv1/math3"
)

type Camera struct {
	pos                math3.Vector
	forward, up, right math3.Vector
}

func newCamera(pos math3.Point, lookAt math3.Vector) Camera {
	down := math3.Vec(0.0, -1.0, 0.0)

	forward := lookAt.Sub(pos).Norm()
	right := forward.Cross(down).Norm().ScaleAll(1.5)
	up := forward.Cross(right).Norm().ScaleAll(1.5)
	return Camera{
		pos:     pos,
		forward: forward,
		right:   right,
		up:      up,
	}
}

type Scene struct {
	width, height int

	camera Camera
	things []Thing
	lights []Light
}
