package main

import (
	"fmt"
	"image/color"
	"strconv"
	"strings"
)

// rgbaFromHexString parses the given string into RGBA.
func rgbaFromHexString(str string) (color.RGBA, error) {
	colStr := strings.TrimPrefix(strings.ToLower(str), "0x")
	c, err := strconv.ParseUint(colStr, 16, 32)
	if err != nil {
		return color.RGBA{}, fmt.Errorf("parse uint: %w", err)
	}
	return rgba(uint32(c)), nil
}

// rgba handles the bitshift magic to parse RGBA.
// The magic is needed to support bot 0xFFFFFF and 0xFFFFFFFF.
func rgba(rgba uint32) color.RGBA {
	var cc color.RGBA

	if rgba > 0xFFFFFF {
		cc.A = 8
	}
	cc.R = uint8((rgba >> (16 + cc.A)) & 0xFF)
	cc.G = uint8((rgba >> (8 + cc.A)) & 0xFF)
	cc.B = uint8((rgba >> cc.A) & 0xFF)
	if cc.A != 0 {
		cc.A = uint8(rgba & 0xFF)
	}

	return cc
}
