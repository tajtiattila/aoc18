package constellation

import (
	"fmt"
	"sort"
	"testing"
)

func TestTreeNeighbots(t *testing.T) {
	type nbor struct {
		p Point
		q []Point
	}

	type test struct {
		points []Point
		nbor   []nbor
	}

	tests := []test{
		{
			points: []Point{
				Pt(0, 0, 0, 0),
				Pt(3, 0, 0, 0),
				Pt(0, 3, 0, 0),
				Pt(0, 0, 3, 0),
				Pt(0, 0, 0, 3),
				Pt(0, 0, 0, 6),
				Pt(9, 0, 0, 0),
				Pt(12, 0, 0, 0),
			},
			nbor: []nbor{
				{
					p: Pt(0, 0, 0, 0),
					q: []Point{
						Pt(3, 0, 0, 0),
						Pt(0, 3, 0, 0),
						Pt(0, 0, 3, 0),
						Pt(0, 0, 0, 3),
					},
				},
				{
					p: Pt(0, 0, 0, 3),
					q: []Point{
						Pt(0, 0, 0, 0),
						Pt(0, 0, 0, 6),
					},
				},
			},
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprint(i+1), func(t *testing.T) {
			sortPoints(tt.points)

			tree := NewTree(tt.points)
			t.Logf("%+v", tree)

			for i, nt := range tt.nbor {
				var got []Point
				tree.WalkNeighbors(nt.p, 3, func(e Elem) {
					if e.P != nt.p {
						got = append(got, e.P)
					}
				})

				sortPoints(nt.q)
				sortPoints(got)

				if !samePoints(got, nt.q) {
					t.Fatalf("%d: got %v, want %v", i, got, nt.q)
				}
			}
		})
	}
}

func samePoints(a, b []Point) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}

func sortPoints(p []Point) {
	sort.Slice(p, func(i, j int) bool {
		for axis := 0; axis < Dim; axis++ {
			if p[i][axis] != p[j][axis] {
				return p[i][axis] < p[j][axis]
			}
		}
		return false
	})
}
