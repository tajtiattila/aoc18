package nanobot

import (
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
		m := Mpt(p.x, p.y, p.z)
		if !m.Valid() {
			t.Fatalf("point %d,%d,%d (%d,%d,%d) is invalid",
				p.x, p.y, p.z, m.Xm, m.Ym, m.Zm)
		}
		x, y, z := m.Coords()
		if x != p.x || y != p.y || z != p.z {
			t.Fatalf("got %d,%d,%d (%d,%d,%d); want %d,%d,%d",
				x, y, z, m.Xm, m.Ym, m.Zm, p.x, p.y, p.z)
		}
	}
}

func TestMBox(t *testing.T) {
	for r := 1; r < 10; r++ {

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

		var got []point
		bb := Equidist(0, 0, 0, r)
		for xm := bb.Min.Xm; xm < bb.Max.Xm; xm++ {
			for ym := bb.Min.Ym; ym < bb.Max.Ym; ym++ {
				for zm := bb.Min.Zm; zm < bb.Max.Zm; zm++ {
					m := MPoint{Xm: xm, Ym: ym, Zm: zm}
					if m.Valid() {
						x, y, z := m.Coords()
						got = append(got, point{x: x, y: y, z: z})
					}
				}
			}
		}

		sort.Slice(got, func(i, j int) bool {
			if got[i].x != got[j].x {
				return got[i].x < got[j].x
			}
			if got[i].y != got[j].y {
				return got[i].y < got[j].y
			}
			return got[i].z < got[j].z
		})

		if len(got) != len(want) {
			t.Fatalf("with r=%d got len %d; want %d", r, len(got), len(want))
		}
	}
}
