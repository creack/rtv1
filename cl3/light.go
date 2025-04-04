package main

// Light represents a point light source
type Light struct {
	Position  Vector3
	Color     Color
	Intensity float64
}

// CreateLights creates the default set of lights for the scene
func CreateLights() []Light {
	return []Light{
		{
			Position:  Vector3{-2, 2, 0},
			Color:     Color{1, 1, 1},
			Intensity: 5.0,
		},
		{
			Position:  Vector3{2, 1, 0},
			Color:     Color{0.8, 0.8, 1},
			Intensity: 3.0,
		},
		// Add a light to better illuminate the cylinder and cone
		{
			Position:  Vector3{0, 0, 1.5},
			Color:     Color{1, 1, 0.9},
			Intensity: 2.0,
		},
	}
}
