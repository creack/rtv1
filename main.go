// Package main is the entry point of the program.
package main

import (
	_ "image/png"
	"log"
	"os"
)

const (
	initialScreenWidth  = 800
	initialScreenHeight = 600
)

const dumpPNG = false // For debug/testing purposes.

// RenderMode enum type.
type RenderMode int

// RenderMode enum values.
const (
	RenderModeGPU RenderMode = iota
	RenderModeCPU
)

func main() {
	s, err := loadScene("")
	if err != nil {
		log.Fatalf("Failed to load scene scene.json: %s.", err)
	}

	g := &Game{
		scene:    s,
		sceneIdx: 0,

		renderMode: RenderModeGPU,
	}

	// TODO: Document the fake shader aspect.
	// WHen the fake shader mode is enabled, set the render to be CPU.
	if os.Getenv("FAKE_SHADER") == "1" {
		g.renderMode = RenderModeCPU
	}

	// If we are in GPU render mode, compile the shader in the background.
	// Wait a little for the window to be created.
	if g.renderMode == RenderModeGPU {
		g.shader = compileShader(g.scene)
	}

	g.run()
}
