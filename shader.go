package main

import (
	"embed"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

var (
	//go:embed k_*.go
	shaderGo embed.FS
	//go:embed k_*.kage
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
		dumpCompiledShader(str)
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

func dumpCompiledShader(str string) {
	dumpShader := os.Getenv("DUMP_SHADER")
	// If disabled, nothing to do.
	if dumpShader == "0" || dumpShader == "" {
		return
	}

	// If in "2" mode, dump to stderr with line numbers.
	if dumpShader == "2" {
		for i, line := range strings.Split(str, "\n") {
			fmt.Fprintf(os.Stderr, "[%d] %s\n", i, line)
		}
		return
	}

	// Otherwise, dump to a file, trimed down.
	buf := ""
	for _, line := range strings.Split(str, "\n") {
		if !strings.HasPrefix(line, "//kage:") && strings.HasPrefix(strings.TrimSpace(line), "//") {
			continue
		}
		line = strings.Split(line, "//")[0]
		if strings.TrimSpace(line) == "" {
			continue
		}
		buf += line + "\n"
	}

	// Write the file.
	if err := os.WriteFile(dumpShader, []byte(buf), 0644); err != nil {
		log.Printf("Error writing shader file: %s", err)
	}
}
