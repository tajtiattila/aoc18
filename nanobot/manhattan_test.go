package nanobot

import (
	"fmt"
	"sort"
	"testing"
)

type point struct {
	x, y, z int
}

func TestMPoint(t *testing.T) {

	const n = 10

	ch := make(chan point)
	go func() {
		defer close(ch)
		for z := 0; z < n; z++ {
			for y := 0; y < n; y++ {
				for x := 0; x < n; x++ {
					ch <- point{x: x, y: y, z: z}
				}
			}
		}
	}()

	for p := range ch {
		m := MPt(p.x, p.y, p.z)
		if !m.Valid() {
			t.Fatalf("point %d,%d,%d %v is invalid",
				p.x, p.y, p.z, m)
		}
		x, y, z := m.Coords()
		if x != p.x || y != p.y || z != p.z {
			t.Fatalf("got %d,%d,%d %v; want %d,%d,%d",
				x, y, z, m, p.x, p.y, p.z)
		}
	}

	// ensure only one point is valid in the space
	m := make(map[point][]MPoint)
	for x := -n; x < n; x++ {
		for y := -n; y < n; y++ {
			for z := -n; z < n; z++ {
				for w := -n; w < n; w++ {
					p := MPoint{x, y, z, w}
					if p.Valid() {
						x, y, z := p.Coords()
						q := point{x, y, z}
						m[q] = append(m[q], p)
					}
				}
			}
		}
	}

	for k, v := range m {
		if len(v) != 1 {
			t.Fatal(k, MPt(k.x, k.y, k.z), v)
		}
	}
}

func TestMBox(t *testing.T) {
	for r := 1; r < 12; r++ {

		var want []point
		for x := -r; x <= r; x++ {
			for y := -r; y <= r; y++ {
				for z := -r; z <= r; z++ {
					dx, dy, dz := x, y, z
					if dx < 0 {
						dx = -dx
					}
					if dy < 0 {
						dy = -dy
					}
					if dz < 0 {
						dz = -dz
					}
					if dx+dy+dz <= r {
						want = append(want, point{x: x, y: y, z: z})
					}
				}
			}
		}

		bb := Equidist(0, 0, 0, r)

		lw := len(want)
		np := bb.NumPoints()

		t.Logf("raduis %d -> %d points", r, lw)
		t.Log(np, np/lw, np%lw)

		got := mboxCoords(t, bb)

		if len(got) != len(want) {
			t.Fatalf("with r=%d got len %d; want %d\n%v\n%v", r, len(got), len(want), got, want)
		}
	}
}

func TestMBoxIntersect(t *testing.T) {
	type test struct {
		npoints int
		minp    point
		boxes   []MBox
	}
	tests := []test{
		{
			npoints: 0,
			boxes: []MBox{
				Equidist(10, 12, 12, 2),
				Equidist(12, 14, 12, 2),
				Equidist(16, 12, 12, 4),
				Equidist(14, 14, 14, 6),
				Equidist(50, 50, 50, 200),
				Equidist(10, 10, 10, 5),
			},
		},
		{
			npoints: 1,
			minp:    point{x: 12, y: 12, z: 12},
			boxes: []MBox{
				Equidist(10, 12, 12, 2),
				Equidist(12, 14, 12, 2),
				Equidist(16, 12, 12, 4),
				Equidist(14, 14, 14, 6),
				Equidist(50, 50, 50, 200),
			},
		},
	}

	for tn, tt := range tests {
		t.Run(fmt.Sprintf("%d", tn), func(t *testing.T) {
			sum := tt.boxes[0]
			for i, b := range tt.boxes[1:] {
				sum = sum.Intersect(b)
				if tt.npoints > 0 && sum.Empty() {
					t.Fatalf("box %d empty", i+1)
				}
			}

			got := sum.NumPoints()
			if got != tt.npoints {
				t.Fatalf("got numpoints %d; want %d", got, tt.npoints)
			}

			if got != 1 {
				return
			}

			x, y, z, ok := sum.MinPoint()
			gotp := point{x: x, y: y, z: z}
			if !ok || gotp != tt.minp {
				t.Fatalf("got %v min point %v; want %v", ok, gotp, tt.minp)
			}
		})
	}
}

func mboxCoords(t *testing.T, bb MBox) []point {
	m := make(map[point][]MPoint)
	for xm := bb.Min[0]; xm < bb.Max[0]; xm++ {
		for ym := bb.Min[1]; ym < bb.Max[1]; ym++ {
			for zm := bb.Min[2]; zm < bb.Max[2]; zm++ {
				for wm := bb.Min[3]; wm < bb.Max[3]; wm++ {
					p := MPoint{xm, ym, zm, wm}
					if p.Valid() {
						x, y, z := p.Coords()
						q := point{x: x, y: y, z: z}
						m[q] = append(m[q], p)
					}
				}
			}
		}
	}

	v := make([]point, 0, len(m))
	for k, pv := range m {
		v = append(v, k)
		if len(pv) != 1 {
			t.Fatal(k, MPt(k.x, k.y, k.z), pv)
		}
	}

	sort.Slice(v, func(i, j int) bool {
		if v[i].x != v[j].x {
			return v[i].x < v[j].x
		}
		if v[i].y != v[j].y {
			return v[i].y < v[j].y
		}
		return v[i].z < v[j].z
	})

	return v
}
