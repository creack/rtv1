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
var (
	sceneObjects   ThingsT
	sceneLights    LightsT
	sceneMaterials MaterialsT
)

type sphere struct {
	Center   vec3   `json:"center"`
	Radius   float  `json:"radius"`
	Material string `json:"material"`
}

func (s sphere) mat4() mat4 { return newSphere(s.Center, s.Radius, materialTypeIndex[s.Material]) }

func (s sphere) marshalConstructor() string {
	return fmt.Sprintf("newSphere(%s, %f, %d)", s.Center.marshalConstructor(), s.Radius, materialTypeIndex[s.Material])
}

//nolint:unparam // Keeping for reference.
func getSphere(in mat4) (center vec3, radius, radius2 float, materialIdx int) {
	return in[1].vec3, in[0].z, in[0].w, int(in[0].y)
}

type plane struct {
	Center         vec3   `json:"center"`
	Normal         vec3   `json:"normal"`
	IsCheckerboard bool   `json:"is_checkerboard"`
	CheckerSize    float  `json:"checker_size"`
	Material       string `json:"material"`
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

func getPlane(in mat4) (center, normal vec3, isCheckerBoard bool, checkerSize float, materialIdx int) {
	return in[1].vec3, in[2].vec3, in[0].z != 0.0, in[0].w, int(in[0].y)
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

func getLight(in mat4) (center vec3, color vec4, intensity float) {
	return in[0].vec3, in[2], in[0].w
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

func getMaterial(materials MaterialsT, idx int) (mType float, color vec4, ambient, diffuse, specular, specularPower, reflectiveIndex float) {
	m := materials[idx]
	mType = m[0].x
	color = m[1]
	ambient = m[0].y
	diffuse = m[0].z
	specular = m[0].w
	specularPower = m[1].x
	reflectiveIndex = m[1].y
	return mType, color, ambient, diffuse, specular, specularPower, reflectiveIndex
}

type camera struct {
	Origin vec3 `json:"origin"`
	LookAt vec3 `json:"lookAt"`
}

func getCameraComponents(in mat4) (forward, right, up vec3) {
	return in[0].vec3, in[1].vec3, in[2].vec3
}
