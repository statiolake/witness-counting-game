package geom

import "testing"

func c(a, b float64) Coord {
	return NewCoord(a, b)
}

func s(a, b Coord) Segment {
	return Segment{a, b}
}

func TestCCW(t *testing.T) {
	t.Run("CCW", func(t *testing.T) {
		testcases := []struct {
			a, b, c  Coord
			expected int
		}{
			{c(1, 0), c(1, 1), c(0, 1), 1},
			{c(0, 1), c(1, 1), c(1, 0), -1},
			{c(0, 0), c(1, 1), c(2, 2), 0},
		}

		for _, tc := range testcases {
			if ccw := CCW(tc.a, tc.b, tc.c); ccw != tc.expected {
				t.Fatalf(
					"Wrong CCW: %v, %v, %v: expected %v but %v",
					tc.a, tc.b, tc.c, tc.expected, ccw,
				)
			}
		}
	})
}

func TestSegmentCrosses(t *testing.T) {
	t.Run("SegmentCrosses", func(t *testing.T) {
		testcases := []struct {
			a, b     Segment
			expected bool
		}{
			{s(c(1, 0), c(1, 1)), s(c(0, -1), c(2, -1)), false},
			{s(c(1, 0), c(1, 1)), s(c(0, 0), c(2, 0)), false},
			{s(c(1, -1), c(1, 1)), s(c(0, 0), c(2, 0)), true},
		}

		for _, tc := range testcases {
			if crosses := tc.a.Crosses(tc.b); crosses != tc.expected {
				t.Fatalf(
					"Wrong cross: %v and %v: expected %v but %v",
					tc.a, tc.b, tc.expected, crosses,
				)
			}
		}
	})
}
