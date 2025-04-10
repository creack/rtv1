package main

import (
	_ "embed"
	"fmt"
	"strconv"
	"strings"
)

func preprocess(s scene, files ...[]byte) string {
	// Remove the "package main" line from the secondary files.
	str := string(files[0])
	for _, elem := range files[1:] {
		stripped := strings.Replace(string(elem), "package main", "", 1)
		str += stripped
	}

	// Inject the scene objects.
	for _, elem := range []struct {
		k string
		f func() string
	}{
		{"things", s.marshalInjectThings},
		{"lights", s.marshalInjectLights},
		{"materials", s.marshalInjectMaterials},
		{"ambientLight", func() string { return `ambientLight := ` + s.AmbientLight.marshalConstructor() }},
	} {
		str = strings.ReplaceAll(str, "//scene:"+elem.k, elem.f())
	}

	// Replace the custom types with their underlying equivalents.
	for _, elem := range []struct {
		CustomType string
		ItemType   string
		ArraySize  int
	}{
		{"ThingsT", "mat4", len(s.Objects)},
		{"LightsT", "mat4", len(s.Lights)},
		{"MaterialsT", "mat4", len(s.Materials)},
	} {
		underlying := fmt.Sprintf("[%d]"+elem.ItemType, elem.ArraySize)
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
	for i := range maxDepth {
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
