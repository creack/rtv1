package main

import (
	"encoding/json"
	"fmt"
)

// This file is a wrapper for Kage. It mirror the shaderlib_rtv1.kage file to allow
// for the shader to be compile in Go.
// In this part, we have the RTv1 specific functions and types.

type ThingsT []mat4

type LightsT []mat4

type MaterialsT []mat4

// Globals used to pass the scene to the Fragment function.
// In shader mode, it the constructors get injected.
// Not using Uniform as it doesn't support complex types nor arrays.
var (
	sceneObjects   ThingsT
	sceneLights    LightsT
	sceneMaterials MaterialsT
	ambientLight   mat4
)

type sphere struct {
	Center   vec3   `json:"center"`
	Radius   float  `json:"radius"`
	Material string `json:"material"`
}

func (s *sphere) UnmarshalJSON(data []byte) error {
	type alias sphere
	if err := json.Unmarshal(data, (*alias)(s)); err != nil {
		return err
	}
	if s.Radius <= 0 {
		return fmt.Errorf("radius must be greater than 0")
	}
	if s.Center == (vec3{}) {
		return fmt.Errorf("missing 'center'")
	}
	return nil
}

func (s sphere) mat4() mat4 { return newSphere(s.Center, s.Radius, materialTypeIndex[s.Material]) }

func (s sphere) marshalConstructor() string {
	return fmt.Sprintf("newSphere(%s, %f, %d)", s.Center.marshalConstructor(), s.Radius, materialTypeIndex[s.Material])
}

type plane struct {
	Center         vec3   `json:"center"`
	Normal         vec3   `json:"normal"`
	IsCheckerboard bool   `json:"is_checkerboard"`
	CheckerSize    float  `json:"checker_size"`
	Material       string `json:"material"`
}

func (p *plane) UnmarshalJSON(data []byte) error {
	type alias plane
	if err := json.Unmarshal(data, (*alias)(p)); err != nil {
		return err
	}
	if p.IsCheckerboard && p.CheckerSize <= 0 {
		return fmt.Errorf("checker size must be greater than 0")
	}
	if p.Center == (vec3{}) {
		return fmt.Errorf("missing 'center'")
	}
	if p.Normal == (vec3{}) {
		return fmt.Errorf("missing 'normal'")
	}
	return nil
}

func (p plane) mat4() mat4 {
	return newPlane(p.Center, p.Normal, p.IsCheckerboard, p.CheckerSize, materialTypeIndex[p.Material])
}

func (p plane) marshalConstructor() string {
	return fmt.Sprintf("newPlane(%s, %s, %t, %f, %d)",
		p.Center.marshalConstructor(),
		p.Normal.marshalConstructor(),
		p.IsCheckerboard,
		p.CheckerSize,
		materialTypeIndex[p.Material],
	)
}

type cylinder struct {
	Center1  vec3   `json:"center1"`
	Center2  vec3   `json:"center2"`
	Radius   float  `json:"radius"`
	Material string `json:"material"`
}

func (c *cylinder) UnmarshalJSON(data []byte) error {
	type alias cylinder
	if err := json.Unmarshal(data, (*alias)(c)); err != nil {
		return err
	}
	if c.Radius <= 0 {
		return fmt.Errorf("radius must be greater than 0")
	}
	if c.Center1 == (vec3{}) || c.Center2 == (vec3{}) {
		return fmt.Errorf("missing 'center1' or 'center2'")
	}
	return nil
}

func (c cylinder) mat4() mat4 {
	return newCylinder(c.Center1, c.Center2, c.Radius, materialTypeIndex[c.Material])
}

func (c cylinder) marshalConstructor() string {
	return fmt.Sprintf("newCylinder(%s, %s, %f, %d)",
		c.Center1.marshalConstructor(),
		c.Center2.marshalConstructor(),
		c.Radius,
		materialTypeIndex[c.Material],
	)
}

type cone struct {
	Apex     vec3   `json:"apex"`
	Base     vec3   `json:"base"`
	Radius   float  `json:"radius"`
	Material string `json:"material"`
}

func (c *cone) UnmarshalJSON(data []byte) error {
	type alias cone
	if err := json.Unmarshal(data, (*alias)(c)); err != nil {
		return err
	}
	if c.Radius <= 0 {
		return fmt.Errorf("radius must be greater than 0")
	}
	if c.Apex == (vec3{}) || c.Base == (vec3{}) {
		return fmt.Errorf("missing 'apex' or 'base'")
	}
	return nil
}

func (c cone) mat4() mat4 {
	return newCone(c.Base, c.Apex, c.Radius, materialTypeIndex[c.Material])
}

func (c cone) marshalConstructor() string {
	return fmt.Sprintf("newCone(%s, %s, %f, %d)",
		c.Base.marshalConstructor(),
		c.Apex.marshalConstructor(),
		c.Radius,
		materialTypeIndex[c.Material],
	)
}

type light struct {
	Origin    vec3  `json:"origin"`
	Color     vec4  `json:"color"`
	Intensity float `json:"intensity"`
}

func (l light) mat4() mat4 { return newLight(l.Origin, l.Color, l.Intensity) }

func (l light) marshalConstructor() string {
	return fmt.Sprintf("newLight(%s, %s, %f)", l.Origin.marshalConstructor(), l.Color.marshalConstructor(), l.Intensity)
}

// NOTE: Only populated/accessed at the start during the scene loading phase.
//
//	No concurrent access.
var materialTypeIndex = map[string]int{}

type material struct {
	Type            string `json:"type"`
	Color           vec4   `json:"color"`
	Ambient         float  `json:"ambient"`
	Diffuse         float  `json:"diffuse"`
	Specular        float  `json:"specular"`
	SpecularPower   float  `json:"specular_power"`
	ReflectiveIndex float  `json:"reflective_index"`
}

func (m *material) UnmarshalJSON(data []byte) error {
	type alias material
	if err := json.Unmarshal(data, (*alias)(m)); err != nil {
		return err
	}
	if _, ok := materialTypeIndex[m.Type]; ok {
		return fmt.Errorf("duplicate material type: %q", m.Type)
	}
	materialTypeIndex[m.Type] = len(materialTypeIndex)
	return nil
}

func (m material) mat4() mat4 {
	return newMaterial(materialTypeIndex[m.Type], m.Color, m.Ambient, m.Diffuse, m.Specular, m.SpecularPower, m.ReflectiveIndex)
}

func (m material) marshalConstructor() string {
	return fmt.Sprintf("newMaterial(%d, %s, %f, %f, %f, %f, %f)",
		materialTypeIndex[m.Type],
		m.Color.marshalConstructor(),
		m.Ambient,
		m.Diffuse,
		m.Specular,
		m.SpecularPower,
		m.ReflectiveIndex,
	)
}

type camera struct {
	Origin vec3 `json:"origin"`
	LookAt vec3 `json:"lookAt"`
}
