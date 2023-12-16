// Package main is the entrypoint.
package main

import (
	"embed"
	"flag"
	"log"
	"runtime"

	"github.com/hajimehoshi/ebiten/v2"
)

//go:embed scenes/*.json
var scenesDir embed.FS

// Game holds the state.
type Game struct {
	width, height int

	sceenName string
	scene     *Scene

	rt  *RayTracer
	img *ebiten.Image
}

func main() {
	g := &Game{
		width:  800,
		height: 800,
	}
	g.rt = &RayTracer{}

	flag.StringVar(&g.sceenName, "s", "scenes/scene1.json", "Scene file path.")
	flag.Parse()
	if err := g.loadScene(g.sceenName); err != nil {
		log.Fatal(err)
	}

	ebiten.SetWindowSize(g.width*2, g.height*2)
	ebiten.SetWindowTitle("RTv1")
	if runtime.GOOS != "js" {
		ebiten.SetFullscreen(true)
	}
	println("Starting")
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
