package main

// Material represents surface properties
type Material struct {
	Color           Color
	Ambient         float64
	Diffuse         float64
	Specular        float64
	SpecularPower   float64
	ReflectiveIndex float64
}

// CreateMaterials creates standard materials for the scene
func CreateMaterials() map[string]Material {
	return map[string]Material{
		"red": {
			Color:           Color{1.0, 0.2, 0.2},
			Ambient:         0.1,
			Diffuse:         0.7,
			Specular:        0.5,
			SpecularPower:   32,
			ReflectiveIndex: 0.3,
		},
		"blue": {
			Color:           Color{0.2, 0.2, 0.8},
			Ambient:         0.1,
			Diffuse:         0.7,
			Specular:        0.5,
			SpecularPower:   32,
			ReflectiveIndex: 0.3,
		},
		"green": {
			Color:           Color{0.2, 0.8, 0.2},
			Ambient:         0.2,
			Diffuse:         1.0,
			Specular:        0.5,
			SpecularPower:   16,
			ReflectiveIndex: 0.2,
		},
		"yellow": {
			Color:           Color{1.0, 0.9, 0.1},
			Ambient:         0.2,
			Diffuse:         1.0,
			Specular:        0.6,
			SpecularPower:   12,
			ReflectiveIndex: 0.3,
		},
		"orange": {
			Color:           Color{1.0, 0.5, 0.0},
			Ambient:         0.2,
			Diffuse:         1.0,
			Specular:        0.6,
			SpecularPower:   12,
			ReflectiveIndex: 0.3,
		},
		"white": {
			Color:           Color{0.9, 0.9, 0.9},
			Ambient:         0.1,
			Diffuse:         0.8,
			Specular:        0.3,
			SpecularPower:   16,
			ReflectiveIndex: 0.2,
		},
	}
}
