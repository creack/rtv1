package main

import (
	"log"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	ScreenWidth  = 1920
	ScreenHeight = 1080
)

func main() {
	// Initialize camera vectors
	UpdateCameraVectors()

	// Set up Ebiten
	ebiten.SetWindowSize(ScreenWidth, ScreenHeight)
	ebiten.SetWindowTitle("RTv1 - Interactive Ray Tracer")
	ebiten.SetTPS(60)
	ebiten.SetRunnableOnUnfocused(true)

	go func() {
		return
		time.Sleep(100e6)
		ebiten.SetWindowPosition(-1920, 0)
		time.Sleep(100e6)
		//ebiten.SetFullscreen(true)
	}()

	// Start the game
	if err := ebiten.RunGameWithOptions(NewGame(), &ebiten.RunGameOptions{
		// InitUnfocused: true,
	}); err != nil {
		log.Fatal(err)
	}
}
