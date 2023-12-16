package main

import "go.creack.net/rtv1/math3"

var scene0 = Scene{
	things: []Thing{
		&nPlane{math3.Vec(0.0, 1.0, 0.0), 0.0, checkerboardSurface},
		&nSphere{math3.Vec(0.0, 1.0, -0.25), 1.0 * 1.0, shinySurface},
		&nSphere{math3.Vec(-1.0, 0.5, 1.5), 0.5 * 0.5, shinySurface},
	},
	lights: []Light{
		{pos: math3.Vec(-2.0, 2.5, 0.0), color: newColor(0.49, 0.07, 0.07)},
		{pos: math3.Vec(1.5, 2.5, 1.5), color: newColor(0.07, 0.07, 0.49)},
		{pos: math3.Vec(1.5, 2.5, -1.5), color: newColor(0.07, 0.49, 0.071)},
		{pos: math3.Vec(0.0, 3.5, 0.0), color: newColor(0.21, 0.21, 0.35)},
	},
	camera: newCamera(math3.Vec(3.0, 2.0, 4.0), math3.Vec(-1.0, 0.5, 0.0)),
}
