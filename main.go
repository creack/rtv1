// Package main is the entrypoint.
package main

import (
	"embed"
	"flag"
	"log"
	"os"
	"runtime"
	"runtime/pprof"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

//go:embed scenes/*.json
var scenesDir embed.FS

// Game holds the state.
type Game struct {
	width, height int

	sceenName string
	scene     *Scene

	rt  RayTracer
	img *ebiten.Image
}

func main() {
	if runtime.GOOS != "js" {
		f, err := os.Create("profile")
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}
	g := &Game{
		width:  800,
		height: 800,
	}
	g.rt = RayTracer{}

	now := time.Now()
	img := g.frame()
	dur := time.Since(now).String()
	println(dur)
	if runtime.GOOS != "js" {
		return
	}
	g.img = ebiten.NewImageFromImage(img)

	flag.StringVar(&g.sceenName, "s", "scenes/scene1.json", "Scene file path.")
	flag.Parse()
	// if err := g.loadScene(g.sceenName); err != nil {
	// 	log.Fatal(err)
	// }

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
