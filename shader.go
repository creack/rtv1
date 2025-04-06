package main

import (
	_ "embed"
	"log"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

var (
	//go:embed shaderlib_builtin.kage
	shaderlibBuiltin []byte
	//go:embed shaderlib_rtv1.kage
	shaderlibRTv1 []byte
	//go:embed rtv1_base.go
	rtv1Base []byte
	//go:embed rtv1_recursive.go
	rtv1Recursive []byte
)

type shader struct {
	data            *ebiten.Shader
	err             error
	compileDuration time.Duration
}

func compileShader() shader {
	s, err, duration := trackTime(func() (*ebiten.Shader, error) {
		// Preprocess the Go code into Kage shader code.
		str := preprocess(shaderlibBuiltin, shaderlibRTv1, rtv1Base, rtv1Recursive)
		// Compile the shader.
		return ebiten.NewShader([]byte(str))
	})
	if err != nil {
		log.Printf("Error compiling shader: %s.", err)
	}
	return shader{
		data:            s,
		err:             err,
		compileDuration: duration,
	}
}
