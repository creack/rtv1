package main

import (
	"fmt"
	"strconv"
	"strings"

	"go.creack.net/rtv1/math2"
)

func parseMap(mapData []byte) ([][]MapPoint, error) {
	// Start by cleaning up the input, removing blank lines and dup spaces.
	//nolint:prealloc // False positive.
	var grid [][]string
	for _, line := range strings.Split(string(mapData), "\n") {
		if line == "" || line[0] == '#' {
			continue
		}
		var gridLine []string
		for _, elem := range strings.Split(line, " ") {
			if elem == "" {
				continue
			}
			gridLine = append(gridLine, elem)
		}
		grid = append(grid, gridLine)
	}

	// Then for each point, parse the height and optional color.
	//nolint:prealloc // False positive.
	var m [][]MapPoint
	for y, line := range grid {
		var points []MapPoint
		for x, elem := range line {
			h, err := strconv.ParseUint(elem, 16, 64)
			if err != nil {
				return nil, fmt.Errorf("invalid height %q for %d/%d: %w", elem, y, x, err)
			}

			p := MapPoint{
				Point:    math2.Pt(x, y),
				wallType: int(h),
			}

			points = append(points, p)
		}
		m = append(m, points)
	}

	if len(m) == 0 {
		return nil, fmt.Errorf("no points")
	}

	return m, nil
}

// MapPoint represents an individual point for the wireframe.
// 3d vector with color.
type MapPoint struct {
	math2.Point
	wallType int
}
