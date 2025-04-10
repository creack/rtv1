//nolint:gochecknoglobals,revive // Uniform variables must be global.
package main

var UniCameraOrigin, UniCameraLookAt vec3

// NOTE: "Time", "Cursor" and "Resolution" are the uniform variables used by Kageland for demos.
var Time float
var Resolution, Cursor vec2
