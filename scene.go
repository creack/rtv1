package main

import (
	"embed"
	"encoding/json"
	"fmt"
)

//go:embed scenes/*.json
var sceneFiles embed.FS

type objects []any

var factory = map[string]func() any{
	"plane": func() any {
		return &plane{}
	},
	"sphere": func() any {
		return &sphere{}
	},
	"cylinder": func() any {
		return &cylinder{}
	},
	"cone": func() any {
		return &cone{}
	},
}

func (objs *objects) UnmarshalJSON(data []byte) error {
	var arr []json.RawMessage
	if err := json.Unmarshal(data, &arr); err != nil {
		return fmt.Errorf("failed to unmarshal objects: %w", err)
	}

	*objs = make(objects, len(arr))
	for i, obj := range arr {
		// Extract the type of the object first.
		var tmp struct {
			Type string `json:"type"`
		}
		if err := json.Unmarshal(obj, &tmp); err != nil {
			return fmt.Errorf("failed to unmarshal object type: %w", err)
		}

		// Then instantiate the object using the factory.
		if _, ok := factory[tmp.Type]; !ok {
			return fmt.Errorf("unknown object type: %q", tmp.Type)
		}
		target := factory[tmp.Type]()

		// Finally unmarshal into the final object.
		if err := json.Unmarshal(obj, target); err != nil {
			return fmt.Errorf("failed to unmarshal object %q: %w", tmp.Type, err)
		}

		(*objs)[i] = target
	}

	return nil
}

type scene struct {
	name         string
	Camera       camera     `json:"camera"`
	Objects      objects    `json:"objects"`
	AmbientLight light      `json:"ambient_light"`
	Lights       []light    `json:"lights"`
	Materials    []material `json:"materials"`
}

func loadScene(fileName string) (scene, error) {
	if fileName == "" { // If filename is empty, lookup the first scene.
		files, err := sceneFiles.ReadDir("scenes")
		if err != nil {
			return scene{}, fmt.Errorf("failed to read scenes directory: %w", err)
		}
		if len(files) == 0 {
			return scene{}, fmt.Errorf("no scene files found")
		}
		fileName = files[0].Name()
	}

	// Load the file content.
	buf, err := sceneFiles.ReadFile("scenes/" + fileName)
	if err != nil {
		return scene{}, fmt.Errorf("failed to read scene.json: %w", err)
	}

	// Reset the material index.
	materialTypeIndex = map[string]int{}

	// Parse the json. Most of the logic is in scene.UnmarshalJSON.
	var s scene
	if err := json.Unmarshal(buf, &s); err != nil {
		return scene{}, fmt.Errorf("failed to unmarshal scene: %w", err)
	}
	s.name = fileName

	// Update the scene global variables to pass them to the Fragment function.
	sceneObjects = make(ThingsT, 0, len(s.Objects))
	for _, elem := range s.Objects {
		sceneObjects = append(sceneObjects, elem.(interface{ mat4() mat4 }).mat4())
	}
	sceneLights = make(LightsT, 0, len(s.Lights))
	for _, elem := range s.Lights {
		sceneLights = append(sceneLights, elem.mat4())
	}
	sceneMaterials = make(MaterialsT, 0, len(s.Materials))
	for _, elem := range s.Materials {
		sceneMaterials = append(sceneMaterials, elem.mat4())
	}
	ambientLightColor = s.AmbientLight.Color

	return s, nil
}

// Generate objects the constructors for the shader to compile.
//
// Example:
//
//	sceneObjects := ThingsT{
//	  newPlane(newVec3(0.000000, 1.000000, 0.000000), 0.000000, newVec4(1.000000, 1.000000, 1.000000, 1.000000)),
//	  newSphere(newVec3(0.000000, 1.000000, -0.250000), 1.000000, newVec4(1.000000, 1.000000, 0.000000, 1.000000)),
//	  newSphere(newVec3(-1.000000, 0.500000, 1.500000), 0.500000, newVec4(1.000000, 0.000000, 0.000000, 1.000000)),
//	}
func (s scene) marshalInjectThings() string {
	injectThings := "sceneObjects := ThingsT{\n"
	for _, obj := range s.Objects {
		injectThings += "\t\t" + obj.(interface{ marshalConstructor() string }).marshalConstructor() + ",\n"
	}
	injectThings += "\t}\n"

	return injectThings
}

// Generate the lights constructors for the shader to compile.
//
// Example:
//
//	sceneLights := LightsT{
//	  newLight(newVec3(-2.000000, 2.500000, 0.000000), newVec4(0.490000, 0.070000, 0.070000, 1.000000)),
//	  newLight(newVec3(1.500000, 2.500000, 1.500000), newVec4(0.070000, 0.070000, 0.490000, 1.000000)),
//	  newLight(newVec3(1.500000, 2.500000, -1.500000), newVec4(0.070000, 0.490000, 0.070000, 1.000000)),
//	  newLight(newVec3(0.000000, 3.500000, -1.500000), newVec4(0.210000, 0.210000, 0.350000, 1.000000)),
//	}
func (s scene) marshalInjectLights() string {
	injectLights := "sceneLights := LightsT{\n"
	for _, obj := range s.Lights {
		injectLights += "\t\t" + obj.marshalConstructor() + ",\n"
	}
	injectLights += "\t}\n"

	return injectLights
}

func (s scene) marshalInjectMaterials() string {
	injectMaterials := "sceneMaterials := MaterialsT{\n"
	for _, obj := range s.Materials {
		injectMaterials += "\t\t" + obj.marshalConstructor() + ",\n"
	}
	injectMaterials += "\t}\n"

	return injectMaterials
}
