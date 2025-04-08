package main

import (
	"bytes"
	"fmt"
	"image"
	"image/png"
	"log"
	"math"
	"os"
	"runtime"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// This file holds the Ebiten specific code.
// Game implements ebiten.Game's interface.
type Game struct {
	shader shader

	renderMode RenderMode
	hideHelp   bool

	scene    scene
	sceneIdx int

	renderedImg image.Image

	width, height int
	forceRedraw   bool
}

// Update implements ebiten.Game's interface.
func (g *Game) Update() error {
	// General controls.
	switch {
	// Exit with ESC.
	case inpututil.IsKeyJustPressed(ebiten.KeyEscape):
		return ebiten.Termination
	// Toggle GPU/CPU mode, reset the rendered image.
	case inpututil.IsKeyJustPressed(ebiten.KeySpace):
		g.renderedImg = nil
		if g.renderMode == RenderModeGPU {
			g.renderMode = RenderModeCPU
		} else {
			g.renderMode = RenderModeGPU
			if g.shader.data == nil {
				g.shader = compileShader(g.scene)
			}
		}
	// Toggle the help message.
	case inpututil.IsKeyJustPressed(ebiten.KeyH):
		g.hideHelp = !g.hideHelp
	// Cycle through the scenes.
	case inpututil.IsKeyJustPressed(ebiten.KeyC):
		scenes, err := sceneFiles.ReadDir("scenes")
		if err != nil {
			return fmt.Errorf("failed to read scenes dir: %w", err)
		}
		g.sceneIdx++

		if g.sceneIdx >= len(scenes) {
			g.sceneIdx = 0
		}

		g.scene, err = loadScene(scenes[g.sceneIdx].Name())
		if err != nil {
			return fmt.Errorf("failed to load scene %s: %w", scenes[g.sceneIdx].Name(), err)
		}
		if g.renderMode == RenderModeGPU {
			g.shader = compileShader(g.scene)
		}

		g.renderedImg = nil
	}

	const rotationSpeed = 0.3
	rotated := false
	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		rotated = true

		_, right, _ := getCameraComponents(newCameraComponents(g.scene.Camera.Origin, g.scene.Camera.LookAt))
		g.scene.Camera.LookAt = add3(g.scene.Camera.LookAt, scale3(right, rotationSpeed))
	}
	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		rotated = true

		_, right, _ := getCameraComponents(newCameraComponents(g.scene.Camera.Origin, g.scene.Camera.LookAt))
		g.scene.Camera.LookAt = add3(g.scene.Camera.LookAt, scale3(right, -rotationSpeed))
	}
	if ebiten.IsKeyPressed(ebiten.KeyUp) {
		rotated = true

		if calculatePitch(g.scene.Camera.Origin, g.scene.Camera.LookAt) < 1.2 {
			g.scene.Camera.LookAt = add3(g.scene.Camera.LookAt, newVec3(0, rotationSpeed, 0))
		}
	}
	if ebiten.IsKeyPressed(ebiten.KeyDown) {
		rotated = true

		if calculatePitch(g.scene.Camera.Origin, g.scene.Camera.LookAt) > -1.2 {
			g.scene.Camera.LookAt = add3(g.scene.Camera.LookAt, newVec3(0, -rotationSpeed, 0))
		}
	}

	const moveSpeed = 0.25

	moved := false
	if ebiten.IsKeyPressed(ebiten.KeyW) {
		moved = true
		forward, _, _ := getCameraComponents(newCameraComponents(g.scene.Camera.Origin, g.scene.Camera.LookAt))
		g.scene.Camera.Origin = add3(g.scene.Camera.Origin, scale3(forward, moveSpeed))
		g.scene.Camera.LookAt = add3(g.scene.Camera.LookAt, scale3(forward, moveSpeed))
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) {
		moved = true
		forward, _, _ := getCameraComponents(newCameraComponents(g.scene.Camera.Origin, g.scene.Camera.LookAt))
		g.scene.Camera.Origin = add3(g.scene.Camera.Origin, scale3(forward, -moveSpeed))
		g.scene.Camera.LookAt = add3(g.scene.Camera.LookAt, scale3(forward, -moveSpeed))
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) {
		moved = true
		_, right, _ := getCameraComponents(newCameraComponents(g.scene.Camera.Origin, g.scene.Camera.LookAt))
		g.scene.Camera.Origin = add3(g.scene.Camera.Origin, scale3(right, moveSpeed))
		g.scene.Camera.LookAt = add3(g.scene.Camera.LookAt, scale3(right, moveSpeed))
	}
	if ebiten.IsKeyPressed(ebiten.KeyA) {
		moved = true
		_, right, _ := getCameraComponents(newCameraComponents(g.scene.Camera.Origin, g.scene.Camera.LookAt))
		g.scene.Camera.Origin = add3(g.scene.Camera.Origin, scale3(right, -moveSpeed))
		g.scene.Camera.LookAt = add3(g.scene.Camera.LookAt, scale3(right, -moveSpeed))
	}
	if ebiten.IsKeyPressed(ebiten.KeyE) {
		moved = true
		g.scene.Camera.Origin = add3(g.scene.Camera.Origin, newVec3(0, moveSpeed, 0))
		g.scene.Camera.LookAt = add3(g.scene.Camera.LookAt, newVec3(0, moveSpeed, 0))
	}
	if ebiten.IsKeyPressed(ebiten.KeyQ) {
		moved = true
		g.scene.Camera.Origin = add3(g.scene.Camera.Origin, newVec3(0, -moveSpeed, 0))
		g.scene.Camera.LookAt = add3(g.scene.Camera.LookAt, newVec3(0, -moveSpeed, 0))
	}

	tainted := moved || rotated
	if g.forceRedraw {
		g.forceRedraw = false
		tainted = true
	}

	if g.renderedImg == nil || tainted {
		op := &ebiten.NewImageOptions{
			Unmanaged: true, // We handle the image ourselves. Needed to render the image from shader.
		}
		width := g.width
		height := g.height
		screen := ebiten.NewImageWithOptions(image.Rect(0, 0, width, height), op)

		g.renderedImg = g.draw(screen)
	}

	return nil
}

// drawCPU draws the scene using the shader code but from the CPU.
// Used to debug/troubleshoot and verify the shader logic.
func (g *Game) drawCPU(screen *ebiten.Image, width, height int) {
	// Populate Uniform variables.
	UniCameraOrigin = g.scene.Camera.Origin
	UniCameraLookAt = g.scene.Camera.LookAt

	UniScreenWidth = width
	UniScreenHeight = height

	// Render.
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	buffer := img.Pix
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			c0 := Fragment(newVec4(float(x), float(y), 0, 0), vec2{}, vec4{})
			off := (y*width + x) * 4
			buffer[off+0] = uint8(min(255, c0.x*255))
			buffer[off+1] = uint8(min(255, c0.y*255))
			buffer[off+2] = uint8(min(255, c0.z*255))
			buffer[off+3] = uint8(255)
		}
	}
	screen.WritePixels(buffer)
}

// drawGPU actually draws the scene using the shader code.
func (g *Game) drawGPU(screen *ebiten.Image, width, height int) {
	shader := g.shader.data
	shaderErr := g.shader.err

	if shaderErr != nil {
		ebitenutil.DebugPrint(screen, "Error Compiling shader:\n"+shaderErr.Error())
		return
	} else if shader == nil {
		ebitenutil.DebugPrint(screen, "Compiling shader...")
		return
	}
	op := &ebiten.DrawRectShaderOptions{}

	op.Uniforms = map[string]any{
		"UniScreenWidth":  width,
		"UniScreenHeight": height,

		"UniCameraOrigin": g.scene.Camera.Origin.uniform(),
		"UniCameraLookAt": g.scene.Camera.LookAt.uniform(),
	}

	screen.DrawRectShader(width, height, g.shader.data, op)
}

// draw the scene.
func (g *Game) draw(screen *ebiten.Image) image.Image {
	width, height := screen.Bounds().Dx(), screen.Bounds().Dy()

	img, _, duration := trackTime(func() (*ebiten.Image, error) {
		if g.renderMode == RenderModeCPU {
			g.drawCPU(screen, width, height)
		} else {
			g.drawGPU(screen, width, height)
		}

		img := ebiten.NewImage(width, height)

		img.DrawImage(screen, &ebiten.DrawImageOptions{})
		return img, nil
	})

	buf := bytes.NewBuffer(nil)
	if dumpPNG {
		name := "output-real.png"
		if g.renderMode == RenderModeCPU {
			name = "output-fake.png"
		}
		if err := png.Encode(buf, img); err != nil {
			panic(err)
		}
		if runtime.GOOS != "js" {
			if err := os.WriteFile(name, buf.Bytes(), 0o600); err != nil {
				panic(err)
			}
		}
		img0, err := png.Decode(buf)
		if err != nil {
			panic(err)
		}
		return img0
	}

	msg := "\n\n\n"
	msg += fmt.Sprintf("shader enabled: %t", g.renderMode == RenderModeGPU)
	if g.renderMode == 0 && g.shader.compileDuration > 0 {
		msg += fmt.Sprintf(", shader compile time: %s\n", g.shader.compileDuration)
	} else {
		msg += "\n"
	}

	if dumpPNG {
		msg += fmt.Sprintf("png size: %vKB\n", math.Round(float64(buf.Len())/1024.*100.)/100.)
	}
	msg += fmt.Sprintf("drawn in: %s\n", duration)
	msg += fmt.Sprintf("camera origin: %s, lookAt: %s, pitch: %0.2f\n", g.scene.Camera.Origin, g.scene.Camera.LookAt, calculatePitch(g.scene.Camera.Origin, g.scene.Camera.LookAt))

	msg += fmt.Sprintf("w/h: %dx%d\n", width, height)

	msg += "\nControls:\n"
	msg += " - WASD: move\n"
	msg += " - QE: up/down\n"
	msg += " - Arrows: look\n"

	if g.renderMode == 0 {
		msg += " - Space: Change render mode to CPU\n"
	} else {
		msg += " - Space: Change render mode to GPU (shader)\n"
	}
	msg += " - H: Hide this help\n"
	msg += " - C: Change scene\n"
	msg += "\n"
	msg += fmt.Sprintf("Current scene: %s\n", g.scene.name)
	if !g.hideHelp {
		ebitenutil.DebugPrint(img, msg)
	}
	return img
}

// Draw implements ebiten.Game's interface.
func (g *Game) Draw(screen *ebiten.Image) {
	if g.renderedImg == nil {
		return
	}
	op := &ebiten.DrawImageOptions{}
	screen.DrawImage(ebiten.NewImageFromImage(g.renderedImg), op)
	msg := fmt.Sprintf("\nFPS: %0.2f\n", ebiten.ActualFPS())
	msg += fmt.Sprintf("TPS: %0.2f\n", ebiten.ActualTPS())
	if !g.hideHelp {
		ebitenutil.DebugPrint(screen, msg)
	}
}

//nolint:gochecknglobals // Expected to be global, set in init.
var isMobile bool

// Layout implements ebiten.Game's interface.
func (g *Game) Layout(outsideWidth, outsideHeight int) (w, h int) {
	if isMobile { // When on mobile, restrict the resolution to avoid the browser crashing.
		outsideWidth = 800
		outsideHeight = 600
	}
	if outsideWidth > 1920 {
		outsideWidth = 1920
	}
	if outsideHeight > 1080 {
		outsideHeight = 1080
	}
	if g.width != outsideWidth {
		g.forceRedraw = true
		g.width = outsideWidth
	}
	if g.height != outsideHeight {
		g.forceRedraw = true
		g.height = outsideHeight
	}
	return outsideWidth, outsideHeight
}

func (g *Game) run() {
	ebiten.SetWindowTitle("RTv1 - Shader - Go")
	ebiten.SetWindowSize(initialScreenWidth, initialScreenHeight)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	if err := ebiten.RunGameWithOptions(g, &ebiten.RunGameOptions{}); err != nil {
		log.Fatal(err)
	}
}
