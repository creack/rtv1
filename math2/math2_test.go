package math2_test

import (
	"math"
	"testing"

	"go.creack.net/rtv1/math2"
)

func TestGetAngle(t *testing.T) {
	t.Parallel()

	// Test all angles in triangles.
	for i, tc := range []struct {
		triangle [3]math2.Point
		angles   [3]math2.Angle
	}{
		{
			// Right triangle with right angle in p3 (0, 0).
			triangle: [3]math2.Point{
				{0, 3},
				{4, 0},
				{0, 0},
			},
			angles: [3]math2.Angle{
				math2.NewDegAngle(180 - 90 - 36.87),
				math2.NewDegAngle(36.87),
				math2.NewDegAngle(90),
			},
		},
		{
			// Right triangle with right angle in p2.
			triangle: [3]math2.Point{
				{4, 3},
				{4, 0},
				{0, 0},
			},
			angles: [3]math2.Angle{
				math2.NewDegAngle(180 - 90 - 36.87),
				math2.NewDegAngle(90),
				math2.NewDegAngle(36.87),
			},
		},
		{
			// Right triangle with right angle in p1.
			triangle: [3]math2.Point{
				{4, 3},
				{4, 0},
				{0, 3},
			},
			angles: [3]math2.Angle{
				math2.NewDegAngle(90),
				math2.NewDegAngle(180 - 90 - 36.87),
				math2.NewDegAngle(36.87),
			},
		},
	} {
		p1, p2, p3 := tc.triangle[0], tc.triangle[1], tc.triangle[2]

		if expect, got := -tc.angles[0], math2.GetAngle(p1, p2, p3); round(expect) != round(got) {
			t.Errorf("[%d] invalid angle between (%.2f/%.2f, %.2f/%.2f) and (%.2f/%.2f, %.2f/%.2f):\nexpect:\t%.2f\ngot: \t%.2f",
				i, p1.X, p1.Y, p2.X, p2.Y, p1.X, p1.Y, p3.X, p3.Y, round(expect.Degrees()), round(got.Degrees()))
		}
		if expect, got := tc.angles[0], math2.GetAngle(p1, p3, p2); round(expect) != round(got) {
			t.Errorf("[%d] invalid angle between (%.2f/%.2f, %.2f/%.2f) and (%.2f/%.2f, %.2f/%.2f):\nexpect:\t%.2f\ngot: \t%.2f",
				i, p1.X, p1.Y, p3.X, p3.Y, p1.X, p1.Y, p2.X, p2.Y, round(expect.Degrees()), round(got.Degrees()))
		}

		if expect, got := tc.angles[1], math2.GetAngle(p2, p1, p3); round(expect) != round(got) {
			t.Errorf("[%d] invalid angle between (%.2f/%.2f, %.2f/%.2f) and (%.2f/%.2f, %.2f/%.2f):\nexpect:\t%.2f\ngot: \t%.2f",
				i, p1.X, p1.Y, p2.X, p2.Y, p1.X, p1.Y, p3.X, p3.Y, round(expect.Degrees()), round(got.Degrees()))
		}
		if expect, got := -tc.angles[1], math2.GetAngle(p2, p3, p1); round(expect) != round(got) {
			t.Errorf("[%d] invalid angle between (%.2f/%.2f, %.2f/%.2f) and (%.2f/%.2f, %.2f/%.2f):\nexpect:\t%.2f\ngot: \t%.2f",
				i, p1.X, p1.Y, p3.X, p3.Y, p1.X, p1.Y, p2.X, p2.Y, round(expect.Degrees()), round(got.Degrees()))
		}

		if expect, got := -tc.angles[2], math2.GetAngle(p3, p1, p2); round(expect) != round(got) {
			t.Errorf("[%d] invalid angle between (%.2f/%.2f, %.2f/%.2f) and (%.2f/%.2f, %.2f/%.2f):\nexpect:\t%.2f\ngot: \t%.2f",
				i, p1.X, p1.Y, p2.X, p2.Y, p1.X, p1.Y, p3.X, p3.Y, round(expect.Degrees()), round(got.Degrees()))
		}
		if expect, got := tc.angles[2], math2.GetAngle(p3, p2, p1); round(expect) != round(got) {
			t.Errorf("[%d] invalid angle between (%.2f/%.2f, %.2f/%.2f) and (%.2f/%.2f, %.2f/%.2f):\nexpect:\t%.2f\ngot: \t%.2f",
				i, p1.X, p1.Y, p3.X, p3.Y, p1.X, p1.Y, p2.X, p2.Y, round(expect.Degrees()), round(got.Degrees()))
		}
	}
}

func TestCoordsFromAngleLen(t *testing.T) {
	t.Parallel()

	for i, tr := range []struct {
		section string
		cases   [][3]math2.Point
	}{
		{
			section: "not right rectangle",
			cases: [][3]math2.Point{
				{
					{1, 1},
					{2, 3},
					{5, 4},
				},
			},
		},
		{
			section: "Bottom Right and Top Righ quandrants",
			cases: [][3]math2.Point{
				{
					{2, 1},
					{2, -3},
					{5, -3},
				},
			},
		},
		{
			section: "Bottom Right and Bottom Left quandrants",
			cases: [][3]math2.Point{
				{
					{2, -1},
					{-2, -1},
					{-2, -5},
				},
			},
		},
		{
			section: "Bottom Left and Top Left quandrants",
			cases: [][3]math2.Point{
				{
					{-2, -1},
					{-2, 1},
					{-4, 1},
				},
			},
		},
		{
			section: "Top Left and Top Right quandrants",
			cases: [][3]math2.Point{
				{
					{-2, 1},
					{-2, 5},
					{3, 5},
				},
			},
		},
		{
			section: "Top Right and Bottom Right and Bottom Left quandrants",
			cases: [][3]math2.Point{
				{
					{2, 1},
					{2, -5},
					{-3, -5},
				},
			},
		},
		{
			section: "Top Right quadrant",
			cases: [][3]math2.Point{
				{
					{0, 3},
					{4, 0},
					{0, 0},
				},
				{
					{4, 3},
					{4, 0},
					{0, 0},
				},
				{
					{5, 4},
					{5, 1},
					{1, 1},
				},
			},
		},
		{
			section: "Top Left quadrant",
			cases: [][3]math2.Point{
				{
					{0, 3},
					{-4, 0},
					{0, 0},
				},
				{
					{-4, 3},
					{-4, 0},
					{0, 0},
				},
				{
					{-5, 4},
					{-5, 1},
					{-1, 1},
				},
			},
		},
		{
			section: "Bottom Left quadrant",
			cases: [][3]math2.Point{
				{
					{0, -3},
					{-4, 0},
					{0, 0},
				},
				{
					{-4, -3},
					{-4, 0},
					{0, 0},
				},
				{
					{-5, -4},
					{-5, -1},
					{-1, -1},
				},
			},
		},
		{
			section: "Bottom Right quadrant",
			cases: [][3]math2.Point{
				{
					{0, -3},
					{4, 0},
					{0, 0},
				},
				{
					{4, -3},
					{4, 0},
					{0, 0},
				},
				{
					{5, -4},
					{5, -1},
					{1, -1},
				},
			},
		},
	} {
		i, tr := i, tr
		t.Run(tr.section, func(t *testing.T) {
			t.Parallel()

			for k, elem := range tr.cases {
				p1, p2, p3 := elem[0], elem[1], elem[2]
				for j, tc := range []struct {
					// Trying to infer pp2 from origin/pp1.
					origin, pp1, pp2 math2.Point
				}{
					{p1, p2, p3},
					{p1, p3, p2},
					{p2, p1, p3},
					{p2, p3, p1},
					{p3, p1, p2},
					{p3, p2, p1},
				} {
					origin := tc.origin
					pp1 := tc.pp1
					pp2 := tc.pp2

					// Get the angle and length from pp1/pp2 and origin.
					angle := math2.GetAngle(origin, pp1, pp2)
					m := origin.Magnitude(pp2)

					// Infer pp2 from origin, pp1, angle and length.
					pfound := math2.CoordinatesFromAngleDist(origin, pp1, angle, m)
					if expect, got := pp2, pfound; round(expect.X) != round(got.X) || round(expect.Y) != round(got.Y) {
						t.Errorf("[%d][%d][%d] Invalid coordinates.\nExpect:\t%v\nGot:\t%v\nData:\t%v\nAngle:\t%v\nLength:\t%v\n",
							k, i, j, expect, got, tc, round(angle.Degrees()), m)
					}
				}
			}
		})
	}
}

func round[T ~float64](in T) float64 {
	fac := math.Pow(10, float64(2))
	return math.Round(float64(in)*fac) / fac
}
