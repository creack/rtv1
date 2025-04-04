package main

import (
	"image/color"
	"math"
)

// Color represents RGB values
type Color struct {
	R, G, B float64
}

// Add returns the sum of two colors
func (c Color) Add(other Color) Color {
	return Color{c.R + other.R, c.G + other.G, c.B + other.B}
}

// Mul returns the product of a color and a scalar
func (c Color) Mul(scalar float64) Color {
	return Color{c.R * scalar, c.G * scalar, c.B * scalar}
}

// MulColor returns the component-wise product of two colors
func (c Color) MulColor(other Color) Color {
	return Color{c.R * other.R, c.G * other.G, c.B * other.B}
}

// Clamp restricts color values to [0,1]
func (c Color) Clamp() Color {
	return Color{
		math.Max(0, math.Min(1, c.R)),
		math.Max(0, math.Min(1, c.G)),
		math.Max(0, math.Min(1, c.B)),
	}
}

// ToRGBA converts color to RGBA values suitable for Ebiten
func (c Color) ToRGBA() color.RGBA {
	c = c.Clamp()
	return color.RGBA{
		R: uint8(c.R * 255),
		G: uint8(c.G * 255),
		B: uint8(c.B * 255),
		A: 255,
	}
}
