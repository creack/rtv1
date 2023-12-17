package main

import (
	"image"
	"image/color"
	"math"

	"go.creack.net/rtv1/math3"
)

func newColor(r, g, b float64) myColor {
	return myColor{
		r: r,
		g: g,
		b: b,
	}
}

func (g *Game) frame() image.Image {
	img := image.NewRGBA64(image.Rect(0, 0, g.width, g.height))

	scene := scene0
	g.rt.render(img, scene, g.width, g.height)

	return img
}

type Ray struct {
	start math3.Point
	dir   math3.Vector
}

type Thing interface {
	intersect(ray Ray) Intersection
	normal(pos math3.Vector) math3.Vector
	surface() Surface
}

type Intersection struct {
	ray   Ray
	dist  float64
	thing Thing
}

type Light struct {
	pos   math3.Vector
	color myColor
}

func reduceLights(
	thing Thing,
	pos, norm, rd math3.Vector,
	lights []Light,
	f func(thing Thing, pos, norm, rd math3.Vector, c myColor, light Light, scene Scene) myColor,
	c2 myColor,
	scene Scene,
) myColor {
	out := c2
	for i := 0; i < len(lights); i++ {
		out = f(thing, pos, norm, rd, out, lights[i], scene)
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

func (rt *RayTracer) intersection(ray Ray, scene Scene) Intersection {
	closest := math.Inf(1)
	var closestIntersection Intersection

	for i := 0; i < len(scene.things); i++ {
		intersection := scene.things[i].intersect(ray)
		if intersection.dist > 0 && intersection.dist < closest {
			closestIntersection = intersection
			closest = intersection.dist
		}
	}

	return closestIntersection
}

func (rt *RayTracer) testRay(ray Ray, scene Scene) float64 {
	intersection := rt.intersection(ray, scene)
	if intersection.dist != 0 {
		return intersection.dist
	}
	return -1
}

func (rt *RayTracer) traceRay(ray Ray, scene Scene, depth int) myColor {
	intersection := rt.intersection(ray, scene)
	if intersection.dist == 0 {
		return backgroundColor
	}
	return rt.shade(intersection, scene, depth)
}

func (rt *RayTracer) shade(intersection Intersection, scene Scene, depth int) myColor {
	d := intersection.ray.dir
	pos := intersection.ray.start.Add(d.ScaleAll(intersection.dist))
	normal := intersection.thing.normal(pos)
	reflectDir := d.Sub(normal.ScaleAll(normal.Dot(d)).ScaleAll(2))
	naturalColor := plusColor(
		backgroundColor,
		rt.getNaturalColor(intersection.thing, pos, normal, reflectDir, scene),
	)
	reflectedColor := grayColor
	if depth < maxDepth {
		reflectedColor = rt.getReflectionColor(
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

func (rt *RayTracer) getReflectionColor(thing Thing, pos, normal, rd math3.Vector, scene Scene, depth int) myColor {
	return scaleColor(
		thing.surface().reflect(pos),
		rt.traceRay(Ray{start: pos, dir: rd}, scene, depth+1),
	)
}

func (rt *RayTracer) addLight(thing Thing, pos, norm, rd math3.Vector, col myColor, light Light, scene Scene) myColor {
	ldis := light.pos.Sub(pos)
	livec := ldis.Norm()
	neatIsect := rt.testRay(Ray{start: pos, dir: livec}, scene)
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

func (rt *RayTracer) getNaturalColor(thing Thing, pos, norm, rd math3.Vector, scene Scene) myColor {
	return reduceLights(thing, pos, norm, rd, scene.lights, rt.addLight, defaultColor, scene)
}

func (rt *RayTracer) render(img *image.RGBA64, scene Scene, screenWidth, screenHeight int) {
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
			c0 := rt.traceRay(
				Ray{start: scene.camera.pos, dir: getPoint(x, y, scene.camera)},
				scene,
				0,
			)
			img.SetRGBA64(x, y, c0.Color())
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

func (p *nPlane) intersect(ray Ray) Intersection {
	denom := p.pos.Dot(ray.dir)
	if denom >= 0 {
		return Intersection{}
	}

	dist := (p.pos.Dot(ray.start) + p.offset) / -denom
	return Intersection{thing: p, ray: ray, dist: dist}
}

func (p *nPlane) normal(_ math3.Vector) math3.Vector {
	return p.pos
}

func (p *nPlane) surface() Surface {
	return p._surface
}

type nSphere struct {
	pos      math3.Vector
	radius2  float64
	_surface Surface
}

func (p *nSphere) intersect(ray Ray) Intersection {
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
		return Intersection{}
	}

	return Intersection{thing: p, ray: ray, dist: dist}
}

func (p *nSphere) normal(pos math3.Vector) math3.Vector {
	return pos.Sub(p.pos).Norm()
}

func (p *nSphere) surface() Surface {
	return p._surface
}

type Checkerboard struct{}

func (c Checkerboard) diffuse(pos math3.Vector) myColor {
	if (int(math.Floor(pos.Z))+int(math.Floor(pos.X)))%2 != 0 {
		return newColor(1, 1, 1)
	}
	return defaultColor
}

func (c Checkerboard) specular(_ math3.Vector) myColor {
	return newColor(1, 1, 1)
}

func (c Checkerboard) reflect(pos math3.Vector) float64 {
	if (int(math.Floor(pos.Z))+int(math.Floor(pos.X)))%2 != 0 {
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

func (mc myColor) Color() color.RGBA64 {
	return color.RGBA64{
		R: uint16(min(1., mc.r) * float64(math.MaxUint16)),
		G: uint16(min(1., mc.g) * float64(math.MaxUint16)),
		B: uint16(min(1., mc.b) * float64(math.MaxUint16)),
		A: math.MaxUint16,
	}
}
