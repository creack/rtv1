package main

import (
	"fmt"
	"runtime"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// Layout implements ebiten.
func (g Game) Layout(_, _ int) (w, h int) {
	return g.width, g.height
}

// Draw implements ebiten.
func (g *Game) Draw(screen *ebiten.Image) {
	// screen.Fill(color.Black)
	img := g.img

	if false {
		ebitenutil.DebugPrint(img, fmt.Sprintf(`TPS: %0.2f, FPS: %0.2f
Resolution: %dx%d
Scene: %s

Controls:
  C: Cycle scenes
`, ebiten.ActualTPS(), ebiten.ActualFPS(), g.width, g.height, g.sceenName))
	}
	// ebitenutil.DebugPrint(img, fmt.Sprintf("drawn in %v", g.dur))
	op := &ebiten.DrawImageOptions{}
	screen.DrawImage(img, op)
}

// Update implements ebiten.
func (g *Game) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) || inpututil.IsKeyJustPressed(ebiten.KeyQ) {
		if runtime.GOOS != "js" {
			return fmt.Errorf("exit")
		}
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyC) {
		entries, err := scenesDir.ReadDir("scenes")
		if err != nil {
			return fmt.Errorf("readDir: %w", err)
		}
		i := -1
		for ii := 0; ii < len(entries); ii++ {
			if entries[ii].Name() == g.sceenName {
				i = ii
				break
			}
		}
		_ = i
		// if err := g.loadScene("scenes/" + entries[(i+1)%len(entries)].Name()); err != nil {
		// 	return fmt.Errorf("loadMap: %w", err)
		// }
	}

	return nil
}
