// Package main is the entry point of the program.
package main

import (
	"bytes"
	_ "embed"
	"fmt"
	"image"
	"image/png"
	_ "image/png"
	"log"
	"math"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

//go:embed default.go
var defaultKage []byte

//go:embed rec.go
var recKage []byte

//go:embed shaderlib.kage
var shaderlib []byte

const (
	screenWidth  = 800
	screenHeight = 600

	dumpPNG = false
)

// Game implements ebiten.Game's interface.
type Game struct {
	shaderMu          sync.RWMutex
	shader            *ebiten.Shader
	shaderErr         error
	shaderCompileTime time.Duration
	forceRedraw       chan struct{}

	renderMode int // 0: shader, 1: cpu.
	hideHelp   bool

	cameraOrigin vec3
	cameraLookAt vec3

	sphereOrigins []vec3

	renderedImg image.Image
}

// Update implements ebiten.Game's interface.
func (g *Game) Update() error {
	tainted := true

	switch {
	case inpututil.IsKeyJustPressed(ebiten.KeySpace):
		g.renderedImg = nil
		if g.renderMode == 0 {
			g.renderMode = 1
		} else {
			g.renderMode = 0
			go g.compileShader()
		}
	case inpututil.IsKeyJustPressed(ebiten.KeyH):
		g.hideHelp = !g.hideHelp

	case ebiten.IsKeyPressed(ebiten.KeyShift) && ebiten.IsKeyPressed(ebiten.KeyW):
		g.cameraLookAt.y += 1.
	case ebiten.IsKeyPressed(ebiten.KeyShift) && ebiten.IsKeyPressed(ebiten.KeyS):
		g.cameraLookAt.y -= 1.
	case ebiten.IsKeyPressed(ebiten.KeyShift) && ebiten.IsKeyPressed(ebiten.KeyA):
		g.cameraLookAt.z += 1.
	case ebiten.IsKeyPressed(ebiten.KeyShift) && ebiten.IsKeyPressed(ebiten.KeyD):
		g.cameraLookAt.z -= 1.
	case ebiten.IsKeyPressed(ebiten.KeyShift) && ebiten.IsKeyPressed(ebiten.KeyQ):
		g.cameraLookAt.x += 1.
	case ebiten.IsKeyPressed(ebiten.KeyShift) && ebiten.IsKeyPressed(ebiten.KeyE):
		g.cameraLookAt.x -= 1.

	case ebiten.IsKeyPressed(ebiten.KeyShift) && ebiten.IsKeyPressed(ebiten.KeyUp):
		g.sphereOrigins[1].y += 1.
	case ebiten.IsKeyPressed(ebiten.KeyShift) && ebiten.IsKeyPressed(ebiten.KeyDown):
		g.sphereOrigins[1].y -= 1.
	case ebiten.IsKeyPressed(ebiten.KeyShift) && ebiten.IsKeyPressed(ebiten.KeyRight):
		g.sphereOrigins[1].x -= 1.
	case ebiten.IsKeyPressed(ebiten.KeyShift) && ebiten.IsKeyPressed(ebiten.KeyLeft):
		g.sphereOrigins[1].x += 1.
	case ebiten.IsKeyPressed(ebiten.KeyShift) && ebiten.IsKeyPressed(ebiten.KeyPageUp):
		g.sphereOrigins[1].z += 1.
	case ebiten.IsKeyPressed(ebiten.KeyShift) && ebiten.IsKeyPressed(ebiten.KeyPageDown):
		g.sphereOrigins[1].z -= 1.

	case ebiten.IsKeyPressed(ebiten.KeyS):
		g.cameraOrigin.z += 1.
		g.cameraLookAt.z += 1.
	case ebiten.IsKeyPressed(ebiten.KeyW):
		g.cameraOrigin.z -= 1.
		g.cameraLookAt.z -= 1.
	case ebiten.IsKeyPressed(ebiten.KeyQ):
		g.cameraOrigin.y += 1.
		g.cameraLookAt.y += 1.
	case ebiten.IsKeyPressed(ebiten.KeyE):
		g.cameraOrigin.y -= 1.
		g.cameraLookAt.y -= 1.
	case ebiten.IsKeyPressed(ebiten.KeyA):
		g.cameraOrigin.x += 1.
		g.cameraLookAt.x += 1.
	case ebiten.IsKeyPressed(ebiten.KeyD):
		g.cameraOrigin.x -= 1.
		g.cameraLookAt.x -= 1.

	case ebiten.IsKeyPressed(ebiten.KeyUp):
		g.sphereOrigins[0].y += 1.
	case ebiten.IsKeyPressed(ebiten.KeyDown):
		g.sphereOrigins[0].y -= 1.
	case ebiten.IsKeyPressed(ebiten.KeyRight):
		g.sphereOrigins[0].x -= 1.
	case ebiten.IsKeyPressed(ebiten.KeyLeft):
		g.sphereOrigins[0].x += 1.
	case ebiten.IsKeyPressed(ebiten.KeyPageUp):
		g.sphereOrigins[0].z += 1.
	case ebiten.IsKeyPressed(ebiten.KeyPageDown):
		g.sphereOrigins[0].z -= 1.

	default:
		tainted = false
	}

	select {
	case <-g.forceRedraw:
		tainted = true
	default:
	}

	if g.renderedImg == nil || tainted {
		op := &ebiten.NewImageOptions{}
		op.Unmanaged = true
		screen := ebiten.NewImageWithOptions(image.Rect(0, 0, screenWidth, screenHeight), op)
		g.renderedImg = g.draw(screen)
	}

	return nil
}

func (g *Game) drawFake(screen *ebiten.Image, width, height int) {
	// Populate Uniform variables.
	UniCameraOrigin = g.cameraOrigin
	UniCameraLookAt = g.cameraLookAt

	UniSphereOrigins1 = g.sphereOrigins[0]
	UniSphereOrigins2 = g.sphereOrigins[1]

	// Render.
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	buffer := img.Pix
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			c0 := Fragment(newVec4(float(x), float(y), 0, 0), vec2{}, vec4{})
			off := (y*width + x) * 4
			buffer[off+0] = uint8(min(255, c0.x*255))
			buffer[off+1] = uint8(min(255, c0.y*255))
			buffer[off+2] = uint8(min(255, c0.z*255))
			buffer[off+3] = uint8(255)
		}
	}
	screen.WritePixels(buffer)
}

func (g *Game) drawReal(screen *ebiten.Image, width, height int) {
	g.shaderMu.RLock()
	shader := g.shader
	shaderErr := g.shaderErr
	g.shaderMu.RUnlock()
	if shaderErr != nil {
		ebitenutil.DebugPrint(screen, "Error Compiling shader:\n"+shaderErr.Error())
		return
	} else if shader == nil {
		ebitenutil.DebugPrint(screen, "Compiling shader...")
		return
	}
	op := &ebiten.DrawRectShaderOptions{}

	uniSphereOrigins := make([][]float32, len(g.sphereOrigins))
	for i := range uniSphereOrigins {
		uniSphereOrigins[i] = g.sphereOrigins[i].uniform()
	}
	op.Uniforms = map[string]any{
		"UniCameraOrigin": g.cameraOrigin.uniform(),
		"UniCameraLookAt": g.cameraLookAt.uniform(),
	}
	for i, elem := range uniSphereOrigins {
		op.Uniforms[fmt.Sprintf("UniSphereOrigins%d", i+1)] = elem
	}
	screen.DrawRectShader(width, height, g.shader, op)
}

// draw the scene.
func (g *Game) draw(screen *ebiten.Image) image.Image {
	width, height := screen.Bounds().Dx(), screen.Bounds().Dy()
	start := time.Now()
	if g.renderMode == 1 {
		g.drawFake(screen, width, height)
	} else {
		g.drawReal(screen, width, height)
	}

	img := ebiten.NewImage(width, height)
	op := &ebiten.DrawImageOptions{}
	img.DrawImage(screen, op)

	buf := bytes.NewBuffer(nil)
	if dumpPNG {
		name := "output-real.png"
		if g.renderMode == 1 {
			name = "output-fake.png"
		}
		if err := png.Encode(buf, img); err != nil {
			panic(err)
		}
		if runtime.GOOS != "js" {
			if err := os.WriteFile(name, buf.Bytes(), 0o600); err != nil {
				panic(err)
			}
		}
		img0, err := png.Decode(buf)
		if err != nil {
			panic(err)
		}
		return img0
	}

	msg := "\n\n\n"
	msg += fmt.Sprintf("shader enabled: %t\n", g.renderMode == 0)
	if g.renderMode == 0 && g.shaderCompileTime > 0 {
		msg += fmt.Sprintf("shader compile time: %s\n", g.shaderCompileTime)
	}
	if dumpPNG {
		msg += fmt.Sprintf("png size: %vKB\n", math.Round(float64(buf.Len())/1024.*100.)/100.)
	}
	msg += fmt.Sprintf("drawn in: %s\n", time.Since(start))
	msg += fmt.Sprintf("camera origin: %v\n", g.cameraOrigin)
	msg += fmt.Sprintf("camera lookAt: %v\n", g.cameraLookAt)
	msg += fmt.Sprintf("sphere0 origin: %v\n", g.sphereOrigins[0])
	msg += fmt.Sprintf("sphere1 origin: %v\n", g.sphereOrigins[1])

	msg += "\nControls:\n"
	msg += " - WASDQE: move camera\n"
	msg += " - Shift WASDQE: change camera direction\n"
	msg += " - Arrows/PgUp/PgDown: move sphere 1\n"
	msg += " - Shift Arrows/Shift Pgup/Shift PgDown: move sphere 2\n"
	if g.renderMode == 0 {
		msg += " - Space: Change render mode to CPU\n"
	} else {
		msg += " - Space: Change render mode to GPU (shader)\n"
	}
	msg += " - H: Hide this help\n"
	if !g.hideHelp {
		ebitenutil.DebugPrint(img, msg)
	}
	return img
}

// Draw implements ebiten.Game's interface.
func (g *Game) Draw(screen *ebiten.Image) {
	if g.renderedImg == nil {
		return
	}
	op := &ebiten.DrawImageOptions{}
	screen.DrawImage(ebiten.NewImageFromImage(g.renderedImg), op)
	msg := fmt.Sprintf("\nFPS: %0.2f\n", ebiten.ActualFPS())
	msg += fmt.Sprintf("TPS: %0.2f\n", ebiten.ActualTPS())
	if !g.hideHelp {
		ebitenutil.DebugPrint(screen, msg)
	}
}

// Layout implements ebiten.Game's interface.
func (g *Game) Layout(outsideWidth, outsideHeight int) (w, h int) {
	_, _ = outsideWidth, outsideHeight
	return screenWidth, screenHeight
}

func preprocess() string {
	mainFile := defaultKage
	secondaryFiles := [][]byte{shaderlib, recKage}

	// Remove the "package main" line from the secondary files.
	str := string(mainFile)
	for _, elem := range secondaryFiles {
		stripped := strings.Replace(string(elem), "package main", "", 1)
		str += stripped
	}

	// Replace the custom types with their underlying equivalents.
	for _, elem := range []struct {
		CustomType string
		ItemType   string
		ArraySize  int
	}{
		{"ThingsT", "mat4", len(ThingsT{})},
		{"LightsT", "mat4", len(LightsT{})},
	} {
		underlying := elem.ItemType
		if elem.ArraySize > 0 {
			underlying = fmt.Sprintf("[%d]"+elem.ItemType, elem.ArraySize)
		}
		str = strings.ReplaceAll(str, elem.CustomType, underlying)
	}

	// Current recursive function block.
	curFct := ""
	fcts := map[string]string{}

	// Current recursive stop condition block.
	curIf := ""
	ifs := map[string]string{}

	// Regular lines, outside any directive blocks.
	regular := ""

	// "Lexer" part. Populate buffers based on tokens.
	lines := strings.Split(str, "\n")
	for i := 0; i+1 < len(lines); i++ {
		line := lines[i]
		if strings.Contains(line, "//rec:") {
			directiveParts := strings.Split(line, "//rec:")
			if len(directiveParts) != 2 {
				panic(fmt.Errorf("invalid '//rec:' declaration: %q", line))
			}
			directive := directiveParts[1]
			parts := strings.Split(directive, ":")
			if len(parts) != 2 {
				panic(fmt.Errorf("invalid '//rec:' directive: %q", parts))
			}
			switch parts[0] {
			case "func":
				curFct = parts[1]
				if !strings.Contains(lines[i+1], fmt.Sprintf("func %s(", parts[1])) {
					panic(fmt.Errorf("invalid '//rec:func', no matching %q func on next line", parts[1]))
				}
				lines[i+1] = strings.ReplaceAll(
					lines[i+1],
					fmt.Sprintf("func %s(", parts[1]),
					fmt.Sprintf("func %s(", parts[1]+"__REC__"),
				)
			case "endfunc":
				if parts[1] != curFct {
					panic(fmt.Errorf("invalid '//rec:endfunc', expected %q, got %q", curFct, parts[1]))
				}
				if curIf != "" {
					panic(fmt.Errorf("invalid '//rec:endfunc' inside '//rec:if'"))
				}
				curFct = ""

			case "if":
				if curFct == "" {
					panic(fmt.Errorf("invalid '//rec:if' outside '//rec:func'"))
				}
				curIf = parts[1]
			case "endif":
				if curFct == "" {
					panic(fmt.Errorf("invalid '//rec:endif' outside '//rec:func'"))
				}
				if parts[1] != curIf {
					panic(fmt.Errorf("invalid '//rec:endif', expected %q, got %q", curIf, parts[1]))
				}
				fcts[curFct] += "//rec:if:" + curFct + "-" + curIf + "\n"
				curIf = ""

				// Forward call. Each depth level will have the depth number as suffix. Expect the first one which is the original name.
			case "call":
				lines[i+1] = strings.ReplaceAll(lines[i+1], parts[1], parts[1]+"__REC__")
				// Actuall recursive call. Each depth level will have the next depth number as suffix.
			case "rec-call":
				lines[i+1] = strings.ReplaceAll(lines[i+1], parts[1], parts[1]+"__REC+1__")
			default:
				panic(fmt.Errorf("unknown '//rec:' directive: %q", parts[0]))
			}
			continue
		}

		if curFct != "" {
			if curIf != "" {
				ifs[curFct+"-"+curIf] += line + "\n"
			} else {
				fcts[curFct] += line + "\n"
			}
		} else {
			regular += line + "\n"
		}
	}

	// Generate the recursive code.

	out := regular
	for i := 0; i < maxDepth; i++ {
		for _, content := range fcts {
			if i < maxDepth-1 {
				for ifName, ifContent := range ifs {
					content = strings.ReplaceAll(content, "//rec:if:"+ifName, ifContent)
				}
			}
			n := ""
			if i > 0 {
				n = strconv.Itoa(i)
			}
			content = strings.ReplaceAll(content, "__REC__", n)
			content = strings.ReplaceAll(content, "__REC+1__", strconv.Itoa(i+1))
			out += content + "\n"
		}
	}

	return out
}

func (g *Game) compileShader() {
	g.shaderMu.Lock()
	defer g.shaderMu.Unlock()

	if g.shader != nil {
		return
	}
	now := time.Now()
	str := preprocess()
	if runtime.GOOS != "js" {
		fmt.Println(str)
	}
	s, err := ebiten.NewShader([]byte(str))
	duration := time.Since(now)
	if err != nil {
		log.Printf("Error compiling shader: %s.", err)
		g.shaderErr = err
	} else {
		g.shaderCompileTime = duration
		g.shader = s
	}
	g.forceRedraw <- struct{}{}
}

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("RTv1 - Shader - Go")
	g := &Game{
		forceRedraw:  make(chan struct{}, 1),
		cameraOrigin: newVec3(0, 10, 10),
		cameraLookAt: newVec3(0, 0, 0),
		sphereOrigins: []vec3{
			newVec3(0, 1, -0.25),
			newVec3(-1.0, 0.5, 1.5),
		},
		renderMode: 0,
	}
	if os.Getenv("FAKE_SHADER") == "1" {
		g.renderMode = 1
	}

	if g.renderMode == 0 {
		go func() {
			time.Sleep(250 * time.Millisecond)
			g.compileShader()
		}()
	}

	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
