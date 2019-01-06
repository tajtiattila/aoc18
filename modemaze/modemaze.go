package modemaze

import (
	"bytes"
	"io"

	"github.com/tajtiattila/aoc18/astar"
)

type Point struct {
	X, Y int
}

func Pt(x, y int) Point { return Point{X: x, Y: y} }

type Tile uint8

const (
	Rocky Tile = iota
	Wet
	Narrow
)

type Map struct {
	Depth  int
	Target Point

	dx, dy int
	p      []Tile // tiles
}

func New(depth, xtarget, ytarget int) *Map {
	dim := xtarget * 2
	if sy := ytarget * 2; sy > dim {
		dim = sy
	}

	m := &Map{
		Depth:  depth,
		Target: Pt(xtarget, ytarget),

		dx: dim,
		dy: dim,
		p:  make([]Tile, dim*dim),
	}

	const (
		y0m = 16807
		x0m = 48271
	)

	p := make([]int, dim*dim)

	// first set geologic indices
	ofs := 1 // (1, 0)
	for x := 1; x < dim; x++ {
		p[ofs] = x * y0m
		if x != ofs {
			panic("hopp")
		}
		ofs++
	}
	ofs = dim // (0, 1)
	for y := 1; y < dim; y++ {
		p[ofs] = y * x0m
		ofs += dim
	}

	// in case xtarget/ytarget == 0
	ofs = xtarget + dim*ytarget
	p[ofs] = 0

	// fill rest
	ofs = dim + 1 // (1, 1)
	for y := 1; y < dim; y++ {
		o := ofs
		for x := 1; x < dim; x++ {
			if x == xtarget && y == ytarget {
				p[o] = 0
			} else {
				e0 := m.erosionLevel(p[o-1])   // x-1,y
				e1 := m.erosionLevel(p[o-dim]) // x,y-1
				p[o] = e0 * e1
			}
			o++
		}
		ofs += dim
	}

	// generate tiles
	for i, gi := range p {
		switch m.erosionLevel(gi) % 3 {
		case 0:
			m.p[i] = Rocky
		case 1:
			m.p[i] = Wet
		case 2:
			m.p[i] = Narrow
		}
	}

	return m
}

func (m *Map) RiskLevel() int {
	risk := 0
	ofs := 0
	for y := 0; y <= m.Target.Y; y++ {
		o := ofs
		for x := 0; x <= m.Target.X; x++ {
			switch m.p[o] {
			case Rocky:
				// pass
			case Wet:
				risk += 1
			case Narrow:
				risk += 2
			}
			o++
		}
		ofs += m.dx
	}
	return risk
}

func (m *Map) Write(w io.Writer, dx, dy int) error {
	if m.dx < dx {
		dx = m.dx
	}
	if m.dy < dy {
		dy = m.dy
	}
	var buf bytes.Buffer
	ofs := 0
	for y := 0; y < dy; y++ {
		o := ofs
		for x := 0; x < dx; x++ {
			if x == 0 && y == 0 {
				buf.WriteRune('M')
			} else if x == m.Target.X && y == m.Target.Y {
				buf.WriteRune('T')
			} else {
				switch m.p[o] {
				case Rocky:
					buf.WriteRune('.')
				case Wet:
					buf.WriteRune('=')
				case Narrow:
					buf.WriteRune('|')
				}
			}
			o++
		}
		ofs += m.dx
		buf.WriteRune('\n')
		if _, err := buf.WriteTo(w); err != nil {
			return err
		}
		buf.Reset()
	}
	return nil
}

const (
	caveSystemModulo = 20183
)

func (m *Map) erosionLevel(geologicIndex int) int {
	return (geologicIndex + m.Depth) % caveSystemModulo
}

const (
	toolNone = iota
	toolClimbingGear
	toolTorch
)

type pathState struct {
	x, y int
	tool uint8
}

func (m *Map) PathDuration() (minutes int) {

	const (
		costStep       = 1
		costToolSwitch = 7
	)

	start := pathState{
		tool: toolTorch,
	}

	add := func(dst *[]astar.State, cost int, p pathState) {
		dx := p.x - m.Target.X
		if dx < 0 {
			dx = -dx
		}

		dy := p.y - m.Target.Y
		if dy < 0 {
			dy = -dy
		}

		estimate := dx + dy

		if p.tool != toolTorch {
			estimate += costToolSwitch
		}

		*dst = append(*dst, astar.State{
			Point:        p,
			Cost:         cost,
			EstimateLeft: estimate,
		})
	}

	_, tc := astar.FindPath(start, func(p0 astar.Point, dst []astar.State) (adjacents []astar.State) {
		p := p0.(pathState)
		vstride := m.dx
		ofs := p.x + p.y*vstride

		tile := m.p[ofs]

		// switch tool
		add(&dst, costToolSwitch, pathState{
			x:    p.x,
			y:    p.y,
			tool: switchTool(tile, p.tool),
		})

		// north
		if p.y > 0 && canEnter(m.p[ofs-vstride], p.tool) {
			add(&dst, costStep, pathState{
				x:    p.x,
				y:    p.y - 1,
				tool: p.tool,
			})
		}

		// south
		if p.y+1 < m.dy && canEnter(m.p[ofs+vstride], p.tool) {
			add(&dst, costStep, pathState{
				x:    p.x,
				y:    p.y + 1,
				tool: p.tool,
			})
		}

		// west
		if p.x > 0 && canEnter(m.p[ofs-1], p.tool) {
			add(&dst, costStep, pathState{
				x:    p.x - 1,
				y:    p.y,
				tool: p.tool,
			})
		}

		// east
		if p.x+1 < m.dx && canEnter(m.p[ofs+1], p.tool) {
			add(&dst, costStep, pathState{
				x:    p.x + 1,
				y:    p.y,
				tool: p.tool,
			})
		}

		return dst
	})

	return tc
}

func switchTool(tile Tile, tool uint8) uint8 {
	switch tile {
	case Rocky:
		if tool == toolClimbingGear {
			return toolTorch
		}
		return toolClimbingGear
	case Wet:
		if tool == toolNone {
			return toolClimbingGear
		}
		return toolNone
	case Narrow:
		if tool == toolNone {
			return toolTorch
		}
		return toolNone
	}
	panic("impossible")
}

func canEnter(tile Tile, tool uint8) bool {
	switch tile {
	case Rocky:
		return tool == toolClimbingGear || tool == toolTorch
	case Wet:
		return tool == toolClimbingGear || tool == toolNone
	case Narrow:
		return tool == toolTorch || tool == toolNone
	}
	panic("impossible")
}
