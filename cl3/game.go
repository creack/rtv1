package main

import (
	"fmt"
	"image/color"
	"sync"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font/basicfont"
)

// Game represents the Ebiten game state
type Game struct {
	scene          Scene
	image          *ebiten.Image
	quality        int // 1 = full resolution, 2 = half resolution, etc.
	rendering      bool
	frameCount     int
	lastRenderTime time.Time
	isMoving       bool
}

// Update is called every frame to update the game state
func (g *Game) Update() error {
	g.frameCount++

	// Handle keyboard input
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		// Force a high-quality render
		if !g.rendering {
			oldQuality := g.quality
			g.quality = 1 // Highest quality
			g.startRender()
			g.quality = oldQuality
		}
	}

	// Quit on Escape key press
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		return ebiten.Termination
	}

	// Camera rotation with arrow keys
	const rotationSpeed = 0.03
	moved := false

	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		Yaw += rotationSpeed
		moved = true
	}
	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		Yaw -= rotationSpeed
		moved = true
	}
	if ebiten.IsKeyPressed(ebiten.KeyUp) {
		Pitch += rotationSpeed
		// Limit pitch to avoid gimbal lock
		if Pitch > 1.5 {
			Pitch = 1.5
		}
		moved = true
	}
	if ebiten.IsKeyPressed(ebiten.KeyDown) {
		Pitch -= rotationSpeed
		// Limit pitch to avoid gimbal lock
		if Pitch < -1.5 {
			Pitch = -1.5
		}
		moved = true
	}

	// Update camera vectors when rotation changes
	if moved {
		UpdateCameraVectors()
	}

	// Camera movement (FPS style)
	moveSpeed := 0.05

	if ebiten.IsKeyPressed(ebiten.KeyW) {
		// Move forward in the direction camera is facing
		CameraPosition = CameraPosition.Add(CameraForward.Mul(moveSpeed))
		moved = true
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) {
		// Move backward in the direction camera is facing
		CameraPosition = CameraPosition.Sub(CameraForward.Mul(moveSpeed))
		moved = true
	}
	if ebiten.IsKeyPressed(ebiten.KeyA) {
		// Strafe left (swap controls as requested)
		CameraPosition = CameraPosition.Add(CameraRight.Mul(moveSpeed))
		moved = true
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) {
		// Strafe right (swap controls as requested)
		CameraPosition = CameraPosition.Sub(CameraRight.Mul(moveSpeed))
		moved = true
	}
	if ebiten.IsKeyPressed(ebiten.KeyE) {
		// Move up along world up
		CameraPosition = CameraPosition.Add(Vector3{0, moveSpeed, 0})
		moved = true
	}
	if ebiten.IsKeyPressed(ebiten.KeyQ) {
		// Move down along world up
		CameraPosition = CameraPosition.Sub(Vector3{0, moveSpeed, 0})
		moved = true
	}

	// Track if we're currently moving
	g.isMoving = moved

	// Quality adjustment
	if inpututil.IsKeyJustPressed(ebiten.Key1) {
		g.quality = 1
		if !g.rendering {
			g.startRender()
		}
	} else if inpututil.IsKeyJustPressed(ebiten.Key2) {
		g.quality = 2
		if !g.rendering {
			g.startRender()
		}
	} else if inpututil.IsKeyJustPressed(ebiten.Key3) {
		g.quality = 4
		if !g.rendering {
			g.startRender()
		}
	} else if inpututil.IsKeyJustPressed(ebiten.Key4) {
		g.quality = 8
		if !g.rendering {
			g.startRender()
		}
	}

	// Update camera in scene
	g.scene.Camera.Origin = CameraPosition
	g.scene.Camera.LookAt = CameraPosition.Add(CameraForward)
	g.scene.Camera.Up = CameraUp

	// CPU rendering - live updates

	// If we're moving, use a faster, lower-quality render
	renderQuality := g.quality
	if g.isMoving {
		renderQuality = g.quality // Very low quality for movement
	}

	// If we're not currently rendering and either:
	// 1. We're moving OR
	// 2. We just stopped moving OR
	// 3. We haven't rendered in a while
	now := time.Now()
	shouldRender := g.isMoving ||
		(g.lastRenderTime.Add(100*time.Millisecond).Before(now) &&
			(g.frameCount%10 == 0 || g.image == nil))

	shouldRender = g.image == nil

	if !g.rendering && shouldRender {
		// If we just stopped moving, do a high-quality render
		if !g.isMoving && g.lastRenderTime.Add(300*time.Millisecond).Before(now) {
			renderQuality = 1 // Highest quality
		}

		// Start a new render with the determined quality
		oldQuality := g.quality
		g.quality = renderQuality
		g.startRender()
		g.quality = oldQuality // Restore original quality setting

		g.lastRenderTime = now
	}

	return nil
}

// Draw is called every frame to draw the screen
func (g *Game) Draw(screen *ebiten.Image) {
	// Draw the ray-traced image
	if g.image != nil {
		screen.DrawImage(g.image, &ebiten.DrawImageOptions{})
	}

	// Draw UI text
	var status string

	if g.rendering {
		status = "Rendering..."
	} else if g.isMoving {
		status = "Moving"
	} else {
		status = "Ready"
	}

	info := fmt.Sprintf("FPS: %.1f\nStatus: %s\nQuality: 1/%d\nCamera: (%.2f, %.2f, %.2f)\nControls: WASD=Move, QE=Up/Down, Arrows=Look, 1-4=Quality, ESC=Quit",
		ebiten.ActualFPS(), status, g.quality,
		CameraPosition.X, CameraPosition.Y, CameraPosition.Z)

	text.Draw(screen, info, basicfont.Face7x13, 10, 20, color.White)
}

// Layout defines the game's screen dimensions
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return ScreenWidth, ScreenHeight
}

// startRender initiates a new CPU render
func (g *Game) startRender() {
	if g.rendering {
		return
	}

	g.rendering = true

	// Create a local copy of needed variables
	quality := g.quality
	scene := g.scene
	width := ScreenWidth
	height := ScreenHeight

	// Render in a separate goroutine
	go func() {
		// Create a channel to receive rendered pixels
		type PixelResult struct {
			X, Y  int
			Color Color
		}
		results := make(chan PixelResult, width*height/quality/quality)

		// Render pixels concurrently
		var wg sync.WaitGroup
		for y := 0; y < height; y += quality {
			for x := 0; x < width; x += quality {
				wg.Add(1)
				// go func
				func(px, py int) {
					defer wg.Done()
					u := float64(px) / float64(width-1)
					v := 1.0 - float64(py)/float64(height-1) // Flip Y
					ray := scene.Camera.GetRay(u, v)
					color := scene.TraceRay(ray, scene.MaxDepth, x, y)
					results <- PixelResult{px, py, color}
				}(x, y)
			}
		}

		// Close results channel when all goroutines are done
		go func() {
			wg.Wait()
			close(results)
		}()

		// Create new image
		newImage := ebiten.NewImage(width, height)

		// Collect and draw results
		for result := range results {
			rgba := result.Color.ToRGBA()

			// Fill square for lower quality renders
			for dy := 0; dy < quality; dy++ {
				for dx := 0; dx < quality; dx++ {
					x := result.X + dx
					y := result.Y + dy
					if x < width && y < height {
						newImage.Set(x, y, rgba)
					}
				}
			}
		}

		// Update the image after rendering is complete
		g.image = newImage
		g.rendering = false
	}()
}

// NewGame creates and initializes a new game instance
func NewGame() *Game {
	return &Game{
		scene:          CreateScene(ScreenWidth, ScreenHeight),
		image:          nil, // Will be created on first render
		quality:        1,   // Default quality level
		rendering:      false,
		frameCount:     0,
		lastRenderTime: time.Now().Add(-1 * time.Second), // Force initial render
		isMoving:       false,
	}
}
