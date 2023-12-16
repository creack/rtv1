package main

import (
	"image"
	"image/color"
	"image/draw"
	"math"

	"go.creack.net/rtv1/math3"
)

// BaseObject defines the base fields needed for all objects.
type BaseObject struct {
	position math3.Vector
	color    color.Color
}

// NewBaseObject creates the object.
func NewBaseObject(position math3.Vector, color color.Color) BaseObject {
	return BaseObject{position: position, color: color}
}

// Color returns the Object's color.
func (obj BaseObject) Color() color.Color { return obj.color }

// Plan is the object's implemetation for a Plan.
type Plan struct {
	BaseObject
}

// Intersect calculates the distance between the eye and the Object.
func (p Plan) Intersect(dir, origin math3.Vector) float64 {
	// q := 0.
	// if eye.Z != 0 {
	// 	q = -eye.Z / eye.Z
	// }

	switch origin = origin.Sub(p.position); {
	case p.position.X != 0 && dir.X != 0:
		return -origin.X / dir.X
	case p.position.Y != 0 && dir.Y != 0:
		return -origin.Y / dir.Y
	case p.position.Z != 0 && dir.Z != 0:
		return -origin.Z / dir.Z
	}
	return 0
}

// Sphere is the object's implemetation for a Sphere.
type Sphere struct {
	BaseObject
	R int
}

// Intersect calculates the distance between the eye and the Object.
func (obj Sphere) Intersect(dir, origin math3.Vector) float64 {
	oc := origin.Sub(obj.position)
	a := dir.Dot(dir)
	b := 2 * dir.Dot(oc)
	c := oc.Dot(oc) - float64(obj.R*obj.R)
	return SecondDegree(a, b, c)
	// eye = eye.Sub(s.position)
	// // defer eye.Add(s.position)

	// var (
	// 	a = v.X*v.X + v.Y*v.Y + v.Z*v.Z
	// 	b = 2*float64(eye.X)*v.X + float64(eye.Y)*v.Y + float64(eye.Z)*v.Z
	// 	c = eye.X*eye.X + eye.Y*eye.Y + eye.Z*eye.Z - float64(s.R*s.R)
	// )
	// return SecondDegree(a, b, c)
	// orig = orig.Sub(s.position)

	// a := dir.Dot(dir)
	// b := 2 * dir.Dot(orig)
	// c := orig.Dot(orig) - float64(s.R*s.R)

	// return SecondDegree(a, b, c)
}

func newColor(r, g, b float64) myColor {
	return myColor{
		r: r,
		g: g,
		b: g,
	}
}

func (g *Game) frame() image.Image {
	img := image.NewRGBA(image.Rect(0, 0, g.width, g.height))

	scene := Scene{
		things: []Thing{
			nPlane{math3.Vec(0.0, 1.0, 0.0), 0.0, checkerboardSurface},
			nSphere{math3.Vec(0.0, 1.0, -0.25), 1.0 * 1.0, shinySurface},
			nSphere{math3.Vec(-1.0, 0.5, 1.5), 0.5 * 0.5, shinySurface},
		},
		lights: []Light{
			{pos: math3.Vec(-2.0, 2.5, 0.0), color: newColor(0.49, 0.07, 0.07)},
			// {pos: math3.Vec(1.5, 2.5, 1.5), color: newColor(0.07, 0.07, 0.49)},
			// {pos: math3.Vec(1.5, 2.5, -1.5), color: newColor(0.07, 0.49, 0.071)},
			// {pos: math3.Vec(0.0, 3.5, 0.0), color: newColor(0.21, 0.21, 0.35)},
		},
		camera: newCamera(math3.Vec(3.0, 2.0, 4.0), math3.Vec(-1.0, 0.5, 0.0)),
	}
	g.rt.render(img, scene, g.width, g.height)

	return img
}

func (s *Scene) compute(img draw.Image) {
	// for y := 0; y < s.height; y++ {
	// 	for x := 0; x < s.width; x++ {
	// 		img.Set(x, y, s.calc(x, y))
	// 	}
	// }
}

func (s *Scene) calc(x, y int) myColor {
	// var (
	// 	k   float64 = -1
	// 	col myColor = color.Black
	// )

	// ray := math3.Vector{
	// 	X: 100,
	// 	Y: float64(s.width/2 - x),
	// 	Z: float64(s.height/2 - y),
	// }

	// for _, obj := range s.objects {
	// 	// If k == -1, it is our first pass, so if we have a solution, keep it.
	// 	// After that, we check that the solution is smaller than the one we have.
	// 	if tmp := obj.Intersect(ray, s.camera0.position); tmp > 0 && (k == -1 || tmp < k) {
	// 		k = tmp
	// 		col = obj.Color()
	// 	}
	// }
	// return col
	return myColor{}
}

type Ray struct {
	start math3.Point
	dir   math3.Vector
}

type Thing interface {
	intersect(ray Ray) *Intersection
	normal(pos math3.Vector) math3.Vector
	surface() Surface
}

// func (t *Thing) intersect(ray Ray) *Intersection {
// 	return nil
// }

type Intersection struct {
	ray   Ray
	dist  float64
	thing Thing
}

type Light struct {
	pos   math3.Vector
	color myColor
}

func reduceLights(lights []Light, f func(c myColor, light Light) myColor, c2 myColor) myColor {
	out := c2
	for _, elem := range lights {
		out = f(out, elem)
	}
	return out
}

type Surface interface {
	diffuse(pos math3.Vector) myColor
	specular(pos math3.Vector) myColor
	reflect(pos math3.Vector) float64

	roughness() float64
}

const maxDepth = 5

type RayTracer struct{}

func (this *RayTracer) intersection(ray Ray, scene Scene) *Intersection {
	closest := math.Inf(1)
	var closestIntersection *Intersection

	for _, thing := range scene.things {
		intersection := thing.intersect(ray)
		if intersection != nil && intersection.dist < closest {
			closestIntersection = intersection
			closest = intersection.dist
		}
	}

	return closestIntersection
}

func (this *RayTracer) testRay(ray Ray, scene Scene) float64 {
	intersection := this.intersection(ray, scene)
	if intersection != nil {
		return intersection.dist
	}
	return -1
}

func (this *RayTracer) traceRay(ray Ray, scene Scene, depth int) myColor {
	intersection := this.intersection(ray, scene)
	if intersection == nil {
		return backgroundColor
	}
	return this.shade(intersection, scene, depth)
}

func (this *RayTracer) shade(intersection *Intersection, scene Scene, depth int) myColor {
	d := intersection.ray.dir
	pos := intersection.ray.start.Add(d.ScaleAll(intersection.dist))
	normal := intersection.thing.normal(pos)
	reflectDir := d.Sub(normal.ScaleAll(normal.Dot(d)).ScaleAll(2))
	naturalColor := plusColor(
		backgroundColor,
		this.getNaturalColor(intersection.thing, pos, normal, reflectDir, scene),
	)
	reflectedColor := grayColor
	if depth < maxDepth {
		reflectedColor = this.getReflectionColor(
			intersection.thing,
			pos,
			normal,
			reflectDir,
			scene,
			depth,
		)
	}
	return plusColor(naturalColor, reflectedColor)
}

func (this *RayTracer) getReflectionColor(thing Thing, pos, normal, rd math3.Vector, scene Scene, depth int) myColor {
	return scaleColor(
		thing.surface().reflect(pos),
		this.traceRay(Ray{start: pos, dir: rd}, scene, depth+1),
	)
}

func (this *RayTracer) getNaturalColor(thing Thing, pos, norm, rd math3.Vector, scene Scene) myColor {
	addLight := func(col myColor, light Light) myColor {
		ldis := light.pos.Sub(pos)
		livec := ldis.Norm()
		neatIsect := this.testRay(Ray{start: pos, dir: livec}, scene)
		isInShadow := neatIsect != -1 && neatIsect <= ldis.Magnitude()
		if isInShadow {
			return col
		}

		illum := livec.Dot(norm)
		lcolor := defaultColor
		if illum > 0 {
			lcolor = scaleColor(illum, light.color)
		}
		specular := livec.Dot(rd.Norm())
		scolor := defaultColor
		if specular > 0 {
			scolor = scaleColor(
				math.Pow(specular, thing.surface().roughness()),
				light.color,
			)
		}
		return plusColor(
			col,
			plusColor(
				timesColor(thing.surface().diffuse(pos), lcolor),
				timesColor(thing.surface().specular(pos), scolor),
			),
		)
	}
	return reduceLights(scene.lights, addLight, defaultColor)
}

func (this *RayTracer) render(img draw.Image, scene Scene, screenWidth, screenHeight int) {
	getPoint := func(x, y int, camera Camera) math3.Vector {
		recenterX := func(x int) float64 {
			return (float64(x) - float64(screenWidth)/2.0) / 2.0 / float64(screenWidth)
		}
		recenterY := func(y int) float64 {
			return -(float64(y) - float64(screenHeight)/2.0) / 2.0 / float64(screenHeight)
		}
		return camera.forward.Add(camera.right.ScaleAll(recenterX(x)).Add(camera.up.ScaleAll(recenterY(y)))).Norm()
	}

	for y := 0; y < screenHeight; y++ {
		for x := 0; x < screenWidth; x++ {
			c0 := this.traceRay(
				Ray{start: scene.camera.pos, dir: getPoint(x, y, scene.camera)},
				scene,
				0,
			)
			continue
			c := c0.Color()
			_ = c
			img.Set(x, y, c)
		}
	}
}

func scaleColor(k float64, v myColor) myColor {
	return myColor{
		r: v.r * k,
		g: v.g * k,
		b: v.b * k,
	}
}

func plusColor(v1, v2 myColor) myColor {
	return myColor{
		r: v1.r + v2.r,
		g: v1.g + v2.g,
		b: v1.b + v2.b,
	}
}

func timesColor(v1, v2 myColor) myColor {
	return myColor{
		r: v1.r * v2.r,
		g: v1.g * v2.g,
		b: v1.b * v2.b,
	}
}

var (
	grayColor       myColor = newColor(0.3, 0.3, 0.3)
	backgroundColor myColor = newColor(0, 0, 0)
	defaultColor    myColor = newColor(0, 0, 0)
)

type nPlane struct {
	pos      math3.Vector
	offset   float64
	_surface Surface
}

func (p nPlane) intersect(ray Ray) *Intersection {
	denom := p.pos.Dot(ray.dir)
	if denom >= 0 {
		return nil
	}

	dist := (p.pos.Dot(ray.start) + p.offset) / -denom
	return &Intersection{thing: p, ray: ray, dist: dist}
}

func (p nPlane) normal(_ math3.Vector) math3.Vector {
	return p.pos
}

func (p nPlane) surface() Surface {
	return p._surface
}

type nSphere struct {
	pos      math3.Vector
	radius2  float64
	_surface Surface
}

func (p nSphere) intersect(ray Ray) *Intersection {
	eo := p.pos.Sub(ray.start)
	v := eo.Dot(ray.dir)
	dist := 0.

	if v >= 0 {
		disc := p.radius2 - (eo.Dot(eo) - v*v)
		if disc >= 0 {
			dist = v - math.Sqrt(disc)
		}
	}

	if dist == 0 {
		return nil
	}

	return &Intersection{thing: p, ray: ray, dist: dist}
}

func (p nSphere) normal(pos math3.Vector) math3.Vector {
	return pos.Sub(p.pos).Norm()
}

func (p nSphere) surface() Surface {
	return p._surface
}

type Checkerboard struct{}

func (c Checkerboard) diffuse(pos math3.Vector) myColor {
	if int(math.Floor(pos.Z)+math.Floor(pos.X))%2 != 0 {
		return newColor(1, 1, 1)
	}
	return defaultColor
}

func (c Checkerboard) specular(_ math3.Vector) myColor {
	return newColor(1, 1, 1)
}

func (c Checkerboard) reflect(pos math3.Vector) float64 {
	if int(math.Floor(pos.Z)+math.Floor(pos.X))%2 != 0 {
		return 0.1
	}
	return 0.7
}

func (c Checkerboard) roughness() float64 {
	return 150.
}

type Shiny struct{}

func (c Shiny) diffuse(pos math3.Vector) myColor {
	return newColor(1, 1, 1)
}

func (c Shiny) specular(_ math3.Vector) myColor {
	return grayColor
}

func (c Shiny) reflect(pos math3.Vector) float64 {
	return 0.7
}

func (c Shiny) roughness() float64 {
	return 250.
}

var (
	checkerboardSurface Surface = Checkerboard{}
	shinySurface        Surface = Shiny{}
)

type myColor struct {
	r, g, b float64
}

func (mc myColor) Color() color.Color {
	return color.RGBA64{
		R: uint16(math.Floor(min(1., mc.r) * float64(math.MaxUint16))),
		G: uint16(math.Floor(min(1., mc.g) * float64(math.MaxUint16))),
		B: uint16(math.Floor(min(1., mc.b) * float64(math.MaxUint16))),
		A: math.MaxUint16,
	}
}
