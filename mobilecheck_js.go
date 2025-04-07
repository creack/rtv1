package main

import (
	"strings"
	"syscall/js"

	"github.com/hajimehoshi/ebiten/v2"
)

// Mobile browsers (tested on iPad and iPhone, Safari, Chrome, Arc) tend to crash after a few seconds.
// Detect mobile usage so we can artificially limit the resolution to improve stability.
//
//nolint:gochecknoinits // Expected init to set the isMobile variable.
func init() {
	ua := js.Global().Get("navigator").Get("userAgent").String()
	isMobile = strings.Contains(ua, "Android") || strings.Contains(ua, "iPhone") || strings.Contains(ua, "iPad") || strings.Contains(ua, "iPod")
	if isMobile {
		ebiten.SetTPS(30)
	}
}
