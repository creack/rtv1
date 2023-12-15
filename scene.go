package main

import (
	"fmt"
)

func (g *Game) loadScene(sceenName string) error {
	buf, err := scenesDir.ReadFile(sceenName)
	if err != nil {
		return fmt.Errorf("read scene file %q: %w", sceenName, err)
	}
	_ = buf
	g.sceenName = sceenName
	return nil
}
