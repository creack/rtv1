package main

import (
	"fmt"
	"image/color"
	"math"
	"strings"

	"go.creack.net/rtv1/math3"
)

// Object is the Object's interface.
type Object interface {
	Color() color.Color
	Intersect(v, cameraPos math3.Vector) float64
	// Parse(values ObjectConfig) (Object, error)
}

type Camera0 struct {
	position math3.Vector
	rotation math3.Vector

	forward, up, right math3.Vector
}

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
	camera0       Camera0
	objects       []Object

	camera Camera
	things []Thing
	lights []Light
}

func (g *Game) loadScene(sceenName string) error {
	buf, err := scenesDir.ReadFile(sceenName)
	if err != nil {
		return fmt.Errorf("read scene file %q: %w", sceenName, err)
	}
	_ = buf
	g.sceenName = strings.TrimPrefix(sceenName, "scenes/")
	g.scene = &Scene{
		width:  g.width,
		height: g.height,
		camera0: Camera0{
			position: math3.Vec(-300, 0, 100),
			rotation: math3.Vec(0, 0, 0),
		},
		objects: []Object{
			&Plan{
				BaseObject: BaseObject{
					position: math3.Vec(0, 0, 200),
					color:    rgba(0x00FF00),
				},
			},
			// &Plan{
			// 	BaseObject: BaseObject{
			// 		position: math3.Vec(0, 100, 0),
			// 		color:    rgba(0xFF00FF),
			// 	},
			// },
			// &Plan{
			// 	BaseObject: BaseObject{
			// 		position: math3.Vec(0, 0, 50),
			// 		color:    rgba(0xFF0000),
			// 	},
			// },
			// &Plan{
			// 	position: math3.Vec(0, 0, 100),
			// 	color:    rgba(0xFF0000),
			// },
			// &Plan{
			// 	position: math3.Vec(0, 0, -100),
			// 	color:    rgba(0x00FF00),
			// },
			&Sphere{
				BaseObject: BaseObject{
					position: math3.Vec(100, 0, 0),
					color:    rgba(0xFFFF00),
				},
				R: 100,
			},
			&Sphere{
				BaseObject: BaseObject{
					position: math3.Vec(100, 200, -10),
					color:    rgba(0xFF00FF),
				},
				R: 150,
			},
			&Sphere{
				BaseObject: BaseObject{
					position: math3.Vec(0, 150, 50+50),
					color:    rgba(0x00FFFF),
				},
				R: 60,
			},
			// &Sphere{
			// 	BaseObject: BaseObject{
			// 		position: math3.Vec(100, 200, -10),
			// 		color:    rgba(0x0000FF),
			// 	},
			// 	R: 150,
			// },
		},
	}
	return nil
}

// delta is a small helper to determine the delta
// of the given 2nd degree equation.
// (a^2 * x + b * x + c = 0
func delta(a, b, c float64) float64 {
	// b^2 - 4ac
	return b*b - 4*a*c
}

// SecondDegree is a small helper to solve
// the given 2nd degree equation.
// (a * x^2 + b * x + c = 0
// Returns the smallest positive solution.
// Returns 0 when no solution.
func SecondDegree(a, b, c float64) float64 {
	delta := delta(a, b, c)
	// delta negative: no solution
	if delta < 0 {
		return 0
	}
	// Two solution: (-b + sqrt(delta)) / 2a and (-b - sqrl(delta)) / 2a
	var (
		k1 = (-b + math.Sqrt(delta)) / (2 * a)
		k2 = (-b - math.Sqrt(delta)) / (2 * a)
	)
	if k1 > 0 && k1 < k2 {
		return k1
	}
	return k2
}
