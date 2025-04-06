// Package main is the entry point of the program.
package main

import (
	_ "image/png"
	"os"
)

const (
	initialScreenWidth  = 800
	initialScreenHeight = 600

	dumpPNG = false
)

// RenderMode enum type.
type RenderMode int

// RenderMode enum values.
const (
	RenderModeGPU RenderMode = iota
	RenderModeCPU
)

func main() {
	g := &Game{
		cameraOrigin: newVec3(0, 10, 10),
		cameraLookAt: newVec3(0, 0, 0),
		renderMode:   RenderModeGPU,
	}

	// TODO: Document the fake shader aspect.
	// WHen the fake shader mode is enabled, set the render to be CPU.
	if os.Getenv("FAKE_SHADER") == "1" {
		g.renderMode = RenderModeCPU
	}

	// If we are in GPU render mode, compile the shader in the background.
	// Wait a little for the window to be created.
	if g.renderMode == RenderModeGPU {
		g.shader = compileShader()
	}

	g.run()
}
