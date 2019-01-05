package gridregexp

import (
	"io"

	"github.com/tajtiattila/aoc18/pathfind"
)

type coord int16

type point struct {
	x, y coord
}

func pti(x, y int) point {
	return point{x: coord(x), y: coord(y)}
}
func pt(x, y coord) point {
	return point{x: x, y: y}
}

func (p point) next(dir rune) point {
	switch dir {
	case 'N':
		return pt(p.x, p.y-1)
	case 'S':
		return pt(p.x, p.y+1)
	case 'W':
		return pt(p.x-1, p.y)
	case 'E':
		return pt(p.x+1, p.y)
	}
	panic("invalid next direction")
}

func (gr *GridRegexp) Replay(x, y int, f func(dir rune, x, y int)) {
	start := pointset{
		pti(x, y): struct{}{},
	}

	gr.replay(start, f)
}

type pointset map[point]struct{}

func (ps pointset) add(p point) { ps[p] = struct{}{} }

func (gr *GridRegexp) replay(starts pointset, f func(dir rune, x, y int)) pointset {
	if gr == nil {
		return starts
	}

	var nexts pointset

	switch gr.Op {

	case OpLiteral:
		nexts = make(pointset, len(starts))
		for p := range starts {
			for _, r := range gr.Literal {
				p = p.next(r)
				f(r, int(p.x), int(p.y))
			}
			nexts.add(p)
		}

	case OpEmpty:
		nexts = starts

	case OpSelect:
		nexts = make(pointset)
		for _, sub := range gr.Option {
			x := sub.replay(starts, f)
			for p := range x {
				nexts.add(p)
			}
		}
	}

	return gr.Next.replay(nexts, f)
}

func (gr *GridRegexp) Extent() Bounds {
	x, y := 0, 0
	b := Bounds{
		XMin: x,
		XMax: x,
		YMin: y,
		YMax: y,
	}
	gr.Replay(x, y, func(dir rune, x, y int) {
		if x < b.XMin {
			b.XMin = x
		}
		if y < b.YMin {
			b.YMin = y
		}
		if x > b.XMax {
			b.XMax = x
		}
		if y > b.YMax {
			b.YMax = y
		}
	})
	return b
}

func (gr *GridRegexp) Map() *Map {
	x, y := 0, 0
	bb := gr.Extent()

	dx := (bb.XMax - bb.XMin) + 1
	dy := (bb.YMax - bb.YMin) + 1

	m := &Map{
		bb: bb,

		dx: dx,
		dy: dy,
		p:  make([]Tile, dx*dy),
	}

	gr.Replay(x, y, func(dir rune, x, y int) {
		sx, sy := stepfrom(dir, x, y)
		sdoor, door := dirtiles(dir)

		m.orTile(sx, sy, sdoor)
		m.orTile(x, y, door)
	})

	return m
}

// stepfrom returns the source position when stepping
// in direction dir to x, y.
func stepfrom(dir rune, x, y int) (sx, sy int) {
	switch dir {
	case 'N':
		return x, y + 1
	case 'S':
		return x, y - 1
	case 'W':
		return x + 1, y
	case 'E':
		return x - 1, y
	}
	panic("invalid stepfrom direction")
}

func dirtiles(dir rune) (from, to Tile) {
	switch dir {
	case 'N':
		return TileDoorN, TileDoorS
	case 'S':
		return TileDoorS, TileDoorN
	case 'W':
		return TileDoorW, TileDoorE
	case 'E':
		return TileDoorE, TileDoorW
	}
	panic("invalid dirtiles direction")
}

type Map struct {
	bb Bounds

	dx, dy int

	p []Tile
}

func (m *Map) Bounds() Bounds { return m.bb }

func (m *Map) Tile(x, y int) Tile {
	if x < m.bb.XMin || m.bb.XMax < x ||
		y < m.bb.YMin || m.bb.YMax < y {
		return TileEmpty
	}
	return m.p[m.ofs(x, y)]
}

func (m *Map) MaxDoors() int {

	_, maxDist := pathfind.Flood(pt(0, 0), pathfind.Space{
		Adjacent: m.pathfindAdjacents,
	})

	return maxDist
}

func (m *Map) FarRooms(minDist int) int {

	fm, _ := pathfind.Flood(pt(0, 0), pathfind.Space{
		Adjacent: m.pathfindAdjacents,
	})

	n := 0
	for _, dist := range fm {
		if dist >= minDist {
			n++
		}
	}

	return n
}

func (m *Map) pathfindAdjacents(pp pathfind.Place, dst []pathfind.Place) []pathfind.Place {
	p := pp.(point)
	x, y := p.x, p.y
	t := m.p[m.ofs(int(x), int(y))]

	if t&TileDoorN != 0 {
		dst = append(dst, pt(x, y-1))
	}
	if t&TileDoorS != 0 {
		dst = append(dst, pt(x, y+1))
	}
	if t&TileDoorW != 0 {
		dst = append(dst, pt(x-1, y))
	}
	if t&TileDoorE != 0 {
		dst = append(dst, pt(x+1, y))
	}
	return dst
}

// bool canStep reports if one can go from
// room sx, sy to ex, ey in a single step.
func (m *Map) canStep(sx, sy, ex, ey int) bool {
	if sy == ey {
		// EW
		if ex == sx+1 {
			// E
			return m.p[m.ofs(sx, sy)]&TileDoorE != 0
		}
		if ex == sx-1 {
			// W
			return m.p[m.ofs(sx, sy)]&TileDoorW != 0
		}
	}
	if sx == ex {
		// NS
		if ey == sy-1 {
			// N
			return m.p[m.ofs(sx, sy)]&TileDoorN != 0
		}
		if ey == sy+1 {
			// S
			return m.p[m.ofs(sx, sy)]&TileDoorS != 0
		}
	}
	return false
}

func (m *Map) Write(w io.Writer) error {
	// two characters/tile with
	// room for closing columns and newline
	bpr := m.dx*2 + 2
	l0 := make([]byte, bpr)
	l1 := make([]byte, bpr)
	for y := m.bb.YMin; y <= m.bb.YMax; y++ {
		i := 0
		for x := m.bb.XMin; x <= m.bb.XMax; x++ {
			t := m.Tile(x, y)

			l0[i] = '#'

			if t&TileDoorN != 0 {
				l0[i+1] = '-'
			} else {
				l0[i+1] = '#'
			}

			if t&TileDoorW != 0 {
				l1[i] = '|'
			} else {
				l1[i] = '#'
			}

			if t == TileEmpty {
				l1[i+1] = '#' // no room here
			} else {
				if x == 0 && y == 0 {
					l1[i+1] = 'X'
				} else {
					l1[i+1] = '.'
				}
			}

			i += 2
		}
		l0[i] = '#'
		l0[i+1] = '\n'
		l1[i] = '#'
		l1[i+1] = '\n'
		if _, err := w.Write(l0); err != nil {
			return err
		}
		if _, err := w.Write(l1); err != nil {
			return err
		}
	}

	// closing row
	for i := 0; i < bpr-1; i++ {
		l0[i] = '#'
	}
	if _, err := w.Write(l0); err != nil {
		return err
	}
	return nil
}

func (m *Map) ofs(x, y int) int {
	x -= m.bb.XMin
	y -= m.bb.YMin
	return x + y*m.dx
}

func (m *Map) orTile(x, y int, t Tile) {
	m.p[m.ofs(x, y)] |= t
}

type Bounds struct {
	XMin, YMin int // inclusive
	XMax, YMax int // inclusive
}

type Tile uint8

const (
	TileEmpty Tile = 0

	TileDoorN Tile = 0x01
	TileDoorS Tile = 0x02
	TileDoorW Tile = 0x04
	TileDoorE Tile = 0x08
)
