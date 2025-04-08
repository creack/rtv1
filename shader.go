package main

import (
	"embed"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

var (
	//go:embed k_*.go
	shaderGo embed.FS
	//go:embed *.kage
	shaderKage embed.FS
)

type shader struct {
	data            *ebiten.Shader
	err             error
	compileDuration time.Duration
}

func compileShader(s scene) shader {
	var shaderBufs [][]byte

	// Read the embeded files.
	shaderGoEmbed, err := shaderGo.ReadDir(".")
	if err != nil {
		panic(fmt.Errorf("read directory shaderGo %q: %w", err))
	}
	for _, elem := range shaderGoEmbed {
		buf, err := shaderGo.ReadFile(elem.Name())
		if err != nil {
			panic(fmt.Errorf("read file %q: %w", elem.Name(), err))
		}
		shaderBufs = append(shaderBufs, buf)
	}
	shaderKageEmbed, err := shaderKage.ReadDir(".")
	if err != nil {
		panic(fmt.Errorf("read directory shaderKage %q: %w", err))
	}
	for _, elem := range shaderKageEmbed {
		buf, err := shaderKage.ReadFile(elem.Name())
		if err != nil {
			panic(fmt.Errorf("read file %q: %w", elem.Name(), err))
		}
		shaderBufs = append(shaderBufs, buf)
	}

	shaderData, err, duration := trackTime(func() (*ebiten.Shader, error) {
		// Preprocess the Go code into Kage shader code.
		str := preprocess(s, shaderBufs...)
		const dumpShader = false
		if dumpShader {
			for i, line := range strings.Split(str, "\n") {
				fmt.Printf("[%d] %s\n", i, line)
			}
		}
		// Compile the shader.
		return ebiten.NewShader([]byte(str))
	})
	if err != nil {
		log.Printf("Error compiling shader: %s.", err)
	}
	return shader{
		data:            shaderData,
		err:             err,
		compileDuration: duration,
	}
}
