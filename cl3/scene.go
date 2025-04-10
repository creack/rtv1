package main

import (
	"fmt"
	"math"
)

// Scene represents a complete scene with objects and lights
type Scene struct {
	World        HittableList
	Lights       []Light
	Camera       Camera
	AmbientLight Color
	MaxDepth     int
}

// TraceRay computes the color for a ray
func (s Scene) TraceRay(ray Ray, depth int, x, y int) Color {
	if depth <= 0 {
		return Color{0, 0, 0}
	}

	hit, rec := s.World.Hit(ray, 0.001, math.MaxFloat64)
	if !hit {
		return Color{0.1, 0.1, 0.1} // Background color
	}

	// Initialize with ambient light
	result := rec.Material.Color.Mul(rec.Material.Ambient).MulColor(s.AmbientLight)
	//result := rec.Material.Color

	// Process each light
	for _, light := range s.Lights {

		// Calculate light direction and distance
		lightDir := light.Position.Sub(rec.Point)
		lightDistance := lightDir.Length()
		lightDir = lightDir.Normalize()

		// Check for shadows
		shadowRay := Ray{rec.Point, lightDir}
		shadowHit, _ := s.World.Hit(shadowRay, 0.001, lightDistance)

		if shadowHit {
			continue
		}
		// Diffuse lighting
		diffFactor := math.Max(0, rec.Normal.Dot(lightDir))
		diffuse := rec.Material.Color.Mul(rec.Material.Diffuse * diffFactor)

		if first++; first == 2 {
			//first = true
			fmt.Printf("ray: %.3f, x/y: %d,%d\n", ray, x, y)
			fmt.Printf("recPoint: %.4f, lightDir: %.3v, lightDistance: %v, dist: %v\n", rec.Point, lightDir, lightDistance, rec.T)
			fmt.Printf("diffuse: %.3f\n", diffuse)
		}

		// Specular lighting
		viewDir := ray.Direction.Mul(-1).Normalize()
		reflectDir := lightDir.Mul(-1).Reflect(rec.Normal)
		specFactor := math.Pow(math.Max(0, viewDir.Dot(reflectDir)), rec.Material.SpecularPower)
		specular := light.Color.Mul(rec.Material.Specular * specFactor)

		// Combine diffuse and specular
		combined := diffuse.Add(specular)

		// Apply light color and intensity
		combined = combined.MulColor(light.Color).Mul(light.Intensity)

		// Apply distance attenuation (inverse square law)
		attenuation := 1.0 / (lightDistance * lightDistance)
		combined = combined.Mul(attenuation)

		result = result.Add(combined)
	}

	// Calculate reflection
	if rec.Material.ReflectiveIndex > 0 && depth > 0 {
		reflectDir := ray.Direction.Reflect(rec.Normal)
		reflectRay := Ray{rec.Point, reflectDir}
		reflectColor := s.TraceRay(reflectRay, depth-1, x, y)
		result = result.Add(reflectColor.Mul(rec.Material.ReflectiveIndex))
	}

	return result
}

// CreateScene builds the standard scene
func CreateScene(screenWidth, screenHeight int) Scene {
	// Get materials
	materials := CreateMaterials()

	// Create world with objects
	world := HittableList{
		Objects: []Hittable{
			// Cylinder on the left - rotated 45 degrees around X axis
			Cylinder{
				Center1:  Vector3{-0.5, -0.5, -0.7},
				Center2:  Vector3{-0.5, 0.1, -0.3}, // Rotated 45 degrees manually
				Radius:   0.15,
				Material: materials["yellow"],
			},
			// Cone on the right, with top matching cylinder's top
			Cone{
				Apex:     Vector3{0.5, 0.5, -0.7},
				Base:     Vector3{0.5, -0.5, -0.7},
				Radius:   0.15,
				Material: materials["orange"],
			},
			Sphere{
				Center:   Vector3{0, 0, -1},
				Radius:   0.5,
				Material: materials["red"],
			},
			Sphere{
				Center:   Vector3{-1, -0.25, -1.5}, // Halfway through the plane
				Radius:   0.5,
				Material: materials["blue"],
			},
			Sphere{
				Center:   Vector3{1, 0.5, -1.5}, // Raised up in the air
				Radius:   0.5,
				Material: materials["green"],
			},
			Plane{
				Point:          Vector3{0, -0.5, 0},
				Normal:         Vector3{0, 1, 0},
				Material:       materials["white"],
				IsCheckerboard: true,
				CheckerSize:    0.1, // Size of each checker square
			},
		},
	}

	// Create lights
	lights := CreateLights()

	// Create camera
	camera := CreateCamera(screenWidth, screenHeight)
	fmt.Printf(">>>> %v\n", camera.LookAt)

	// Create scene
	return Scene{
		World:        world,
		Lights:       lights,
		Camera:       camera,
		AmbientLight: Color{0.85, 0.1, 0.05}, // Reduced ambient light
		MaxDepth:     3,                      // Reduced for interactive performance
	}
}
