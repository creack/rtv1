package main

//kage:unit pixels

// This file is the main RTv1 logic. It compiles to both Go and Kage shader (after pre-processing).

var UniScreenWidth, UniScreenHeight int

//nolint:gochecknoglobals,revive // Uniform variables must be global.
var UniCameraOrigin, UniCameraLookAt vec3

// getThingType returns the type of the thing.
// By convention, it is stored in the z component of the second column of the mat4.
func getThingType(thing mat4) float {
	return thing[1].z
}

const maxDepth = 5

const (
	SphereType = 1
	PlaneType  = 2
	LightType  = 3
)

// center: s[0].xyz
// radius: s[0].w
// radius2: s[1].w
// roughness: s[1].x
// color: s[2]
// specular: s[3]
func newSphere(center vec3, radius float, col vec4) mat4 {
	return newMat4(
		newVec4(center.x, center.y, center.z, radius),
		newVec4(250., 0, SphereType, radius*radius),
		col,
		newVec4(0.3, 0.3, 0.3, 1),
	)
}

// center: p[0].xyz
// offset: p[0].w
// roughness: p[1].x
// color: p[2]
// specular: p[3]
func newPlane(center vec3, offset float, col vec4) mat4 {
	return newMat4(
		newVec4(center.x, center.y, center.z, offset),
		newVec4(150, 0, PlaneType, 0),
		col,
		newVec4(1, 1, 1, 1),
	)
}

func newLight(center vec3, color vec4) mat4 {
	return newMat4(
		newVec4(center.x, center.y, center.z, 0),
		newVec4(0, 0, LightType, 0),
		color,
		newVec4(0, 0, 0, 0),
	)
}

func newCamera(camStart, camLookAt vec3) mat4 {
	down := newVec3(0.0, -1.0, 0.0)

	forward := normalize3(sub3(camLookAt, camStart))
	right := scale3(normalize3(cross3(forward, down)), 2)
	up := scale3(normalize3(cross3(forward, right)), 1.5)
	return newMat4(
		newVec4(camStart.x, camStart.y, camStart.z, 0),
		newVec4(forward.x, forward.y, forward.z, 0),
		newVec4(right.x, right.y, right.z, 0),
		newVec4(up.x, up.y, up.z, 0),
	)
}

// Fragment is the shader's entry point.
//
//nolint:revive // Unexported return is required by the shader API.
func Fragment(position vec4, _ vec2, _ vec4) vec4 {
	width := UniScreenWidth
	height := UniScreenHeight
	x := int(position.x)
	y := int(position.y)

	things := ThingsT{
		newPlane(newVec3(0, 1.0, 0), 0, newVec4(1, 1, 1, 1)),
		newSphere(newVec3(0, 1, -0.25), 1, newVec4(1, 1, 0, 1)),
		newSphere(newVec3(-1.0, 0.5, 1.5), 0.5, newVec4(1, 0, 0, 1)),
	}

	lights := LightsT{
		// newLight(newVec3(-2.0, 2.5, 0), newVec4(1, 1, 1, 1)),
		newLight(newVec3(-2.0, 2.5, 0), newVec4(0.49, 0.07, 0.07, 1)),
		newLight(newVec3(1.5, 2.5, 1.5), newVec4(0.07, 0.07, 0.49, 1)),
		newLight(newVec3(1.5, 2.5, -1.5), newVec4(0.07, 0.49, 0.07, 1)),
		newLight(newVec3(0, 3.5, 0), newVec4(0.21, 0.21, 0.35, 1)),
	}
	//
	camera := newCamera(UniCameraOrigin, UniCameraLookAt)
	//
	rayDir := initRay(width, height, x, y, camera)
	out := trace(camera, rayDir, lights, things, 0)

	return out
}

func reflectPlane(thing mat4, pos vec3) float {
	_ = thing
	if (int(floor(pos.z))+int(floor(pos.x)))%2 != 0 {
		return 0.1
	}
	return 0.7
}

func diffusePlane(thing mat4, pos vec3) vec4 {
	if (int(floor(pos.z))+int(floor(pos.x)))%2 != 0 {
		return thing[2]
	}
	return newVec4(0, 0, 0, 0)
}

func specularPlane(thing mat4, pos vec3) vec4 {
	_ = pos
	return thing[3]
}

func roughnessPlane(thing mat4, pos vec3) float {
	_ = pos
	return thing[1].x
}

func normalPlane(rayStart, center vec3) vec3 {
	_ = rayStart
	return center
}

func hitPlane(rayStart, rayDir vec3, thing mat4) float {
	ppos, _, _, _ := getSphere(thing) //nolint:dogsled // Expected.
	denom := dot3(ppos, rayDir)
	if denom > 0 {
		return 0
	}
	dist := (dot3(ppos, rayStart) + thing[0].w) / (-denom)
	return dist
}

func reflectSphere(thing mat4, pos vec3) float {
	_ = pos
	_ = thing
	return 0.7
}

func diffuseSphere(thing mat4, pos vec3) vec4 {
	_ = pos
	_, _, _, col := getSphere(thing) //nolint:dogsled // Expected.
	return col
}

func specularSphere(thing mat4, pos vec3) vec4 {
	_ = pos
	return thing[3]
}

func roughnessSphere(thing mat4, pos vec3) float {
	_ = pos
	return thing[1].x
}

func normalSphere(rayStart, center vec3) vec3 {
	return normalize3(sub3(rayStart, center))
}

func hitSphere(rayStart, rayDir vec3, thing mat4) float {
	sphereCenter, _, sphereRadius2, _ := getSphere(thing)

	eo := sub3(sphereCenter, rayStart)
	v := dot3(eo, rayDir)

	disc := sphereRadius2 - (dot3(eo, eo) - v*v)
	if disc >= 0 {
		if dist := v - sqrt(disc); dist >= 0 {
			return dist
		}
	}
	return 0
}

func intersect(rayStart, rayDir vec3, thing mat4) float {
	if t := getThingType(thing); t == SphereType {
		return hitSphere(rayStart, rayDir, thing)
	} else if t == PlaneType {
		return hitPlane(rayStart, rayDir, thing)
	}
	return -1.
}

func intersection(rayStart, rayDir vec3, things ThingsT) (closestThing mat4, closest float) {
	closest = -1.

	for i := 0; i < len(things); i++ {
		if dist := intersect(rayStart, rayDir, things[i]); dist > 0 {
			if closest == -1 || dist < closest {
				closestThing = things[i]
				closest = dist
			}
		}
	}
	if closest == -1. {
		closest = 0.
	}

	return closestThing, closest
}

func testRay(rayStart, rayDir vec3, things ThingsT) float {
	_, dist := intersection(rayStart, rayDir, things)
	if dist != 0 {
		return dist
	}
	return -1
}

func addLight(thing mat4, pos, norm, rd vec3, col vec4, light mat4, things ThingsT) vec4 {
	lightPos, lightColor := getLight(light)

	ldis := sub3(lightPos, pos)
	livec := normalize3(ldis)

	rayStart := pos
	rayDir := livec
	neatIsect := testRay(rayStart, rayDir, things)

	isInShadow := neatIsect != -1 && neatIsect <= length3(ldis)
	if isInShadow {
		return col
	}

	illum := dot3(livec, norm)
	lcolor := newVec4(0, 0, 0, 1) // defaultColor.
	if illum > 0 {
		lcolor = scale4(lightColor, illum)
	}

	specular := dot3(livec, normalize3(rd))
	scolor := newVec4(0, 0, 0, 1) // defaultColor.
	if specular > 0 {
		roughness := 0.
		if t := getThingType(thing); t == SphereType {
			roughness = roughnessSphere(thing, pos)
		} else if t == PlaneType {
			roughness = roughnessPlane(thing, pos)
		} else {
			roughness = -1.
		}
		scolor = scale4(lightColor, pow(specular, roughness))
	}
	var surfaceSpecular, surfaceDiffuse vec4
	if t := getThingType(thing); t == SphereType {
		surfaceSpecular = specularSphere(thing, pos)
		surfaceDiffuse = diffuseSphere(thing, pos)
	} else if t == PlaneType {
		surfaceSpecular = specularPlane(thing, pos)
		surfaceDiffuse = diffusePlane(thing, pos)
	}
	return add4(
		col,
		add4(
			mul4(surfaceDiffuse, lcolor),
			mul4(surfaceSpecular, scolor),
		),
	)
}

func getNaturalColor(thing mat4, pos, norm, rd vec3, lights LightsT, things ThingsT) vec4 {
	defaultColor0 := newVec4(0, 0, 0, 0)

	out := defaultColor0
	for i := 0; i < len(lights); i++ {
		out = addLight(thing, pos, norm, rd, out, lights[i], things)
	}
	return out
}

func initRay(width, height, x, y int, camera mat4) vec3 {
	recenterX := (float(x) - float(width)/2.0) / 2.0 / float(width)
	recenterY := -(float(y) - float(height)/2.0) / 2.0 / float(height)

	_, forward, right, up := getCamera(camera)

	return normalize3(add3(forward, add3(scale3(right, recenterX), scale3(up, recenterY))))
}
