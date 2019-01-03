package main

import (
	"fmt"
	"log"
)

func run06() {
	var data []Point
	for _, s := range input06v {
		var x, y int
		if _, err := fmt.Sscanf(s, "%d, %d", &x, &y); err != nil {
			log.Fatalf("parse %q: %v", s, err)
		}
		data = append(data, Pt(x, y))
	}

	FindNonInfArea(data)
}

type Point struct {
	X, Y int
}

func Pt(x, y int) Point { return Point{X: x, Y: y} }

func FindNonInfArea(points []Point) int {
	ox, oy, dx, dy := analyzeCoords(points)

	// point id is (index in points)+1
	g := newcgrid(dx, dy)

	for i, p := range points {
		id := i + 1
		g.flood(ox+p.X, oy+p.Y, id)
	}

	// find infinite points on the edges
	infs := map[int]struct{}{
		-1: struct{}{},
	}
	checkinf := func(x, y int) {
		e := g.get(x, y)
		infs[e.id] = struct{}{}
	}
	for x := 0; x < g.dx; x++ {
		checkinf(x, 0)
		checkinf(x, g.dy-1)
	}
	for y := 0; y < g.dy; y++ {
		checkinf(0, y)
		checkinf(g.dx-1, y)
	}

	// count area of non-infinite points
	maxn := 0
	for i := range points {
		id := i + 1
		if _, ok := infs[id]; ok {
			// skip infinite
			continue
		}

		n := 0
		for _, e := range g.p {
			if e.id == id {
				n++
			}
		}
		if n > maxn {
			maxn = n
		}
	}
	return maxn
}

// analyzeCoords returns grid offsets and dimensions for points.
func analyzeCoords(points []Point) (ox, oy, dx, dy int) {
	if len(points) == 0 {
		return
	}

	var ix, ax, iy, ay int
	ix = points[0].X
	iy = points[0].Y
	for _, p := range points {
		if p.X < ix {
			ix = p.X
		}
		if p.Y < iy {
			iy = p.Y
		}
		if p.X > ax {
			ax = p.X
		}
		if p.Y > ay {
			ay = p.Y
		}
	}

	tx := ax - ix + 1
	ty := ay - iy + 1
	ox = tx - ix
	oy = ty - iy

	return ox, oy, 3 * tx, 3 * ty
}

type cgrid struct {
	dx, dy int
	p      []cgrident
}

type cgrident struct {
	id   int // closest point id 0: none, -1: tied
	dist int // closest point manhattan distance
}

func newcgrid(dx, dy int) *cgrid {
	return &cgrid{
		dx: dx,
		dy: dy,
		p:  make([]cgrident, dx*dy),
	}
}

func (g *cgrid) ofs(x, y int) int { return x + y*g.dx }

func (g *cgrid) set(x, y int, id, dist int) {
	g.p[g.ofs(x, y)] = cgrident{
		id:   id,
		dist: dist,
	}
}

func (g *cgrid) get(x, y int) cgrident {
	return g.p[g.ofs(x, y)]
}

func (g *cgrid) adj(x, y int, p []Point) []Point {
	if x > 0 {
		p = append(p, Pt(x-1, y))
	}
	if x+1 < g.dx {
		p = append(p, Pt(x+1, y))
	}
	if y > 0 {
		p = append(p, Pt(x, y-1))
	}
	if y+1 < g.dy {
		p = append(p, Pt(x, y+1))
	}
	return p
}

func (g *cgrid) flood(x, y int, id int) {
	dist := 0
	alive := []Point{Pt(x, y)}
	visited := map[Point]struct{}{
		Pt(x, y): struct{}{},
	}

	g.set(x, y, id, dist)

	var next []Point
	var adj []Point

	visit := func(p Point, v int) {
		if _, ok := visited[p]; ok {
			return
		}
		g.set(p.X, p.Y, v, dist)
		visited[p] = struct{}{}
		next = append(next, p)
	}

	for len(alive) != 0 {
		dist++
		for _, q := range alive {
			adj = g.adj(q.X, q.Y, adj[:0])
			for _, p := range adj {
				ent := g.get(p.X, p.Y)
				switch {
				case ent.id == id:
					// already seen

				case ent.id == 0 || ent.dist > dist:
					// empty or farther from other point(s)
					visit(p, id)

				case ent.dist == dist:
					// tie
					visit(p, -1)
				}
			}
		}
		alive, next = next, alive[:0]
	}
}

func (g *cgrid) show() {
	for y := 0; y < g.dy; y++ {
		var p []byte
		for x := 0; x < g.dx; x++ {
			e := g.get(x, y)
			var b byte
			if e.id == -1 {
				b = '.'
			} else if e.id == 0 {
				b = '!'
			} else if e.id < 27 {
				if e.dist == 0 {
					b = 'A' + byte(e.id) - 1
				} else {
					b = 'a' + byte(e.id) - 1
				}
			} else {
				b = '?'
			}
			p = append(p, b)
		}
		fmt.Println(string(p))
	}
}

func FindAreaCloserThan(points []Point, than int) int {
	ox, oy, dx, dy := analyzeCoords(points)

	fx := dx - ox
	fy := dy - oy

	n := 0
	for y := -oy; y < fy; y++ {
		for x := -ox; x < fx; x++ {
			check := Pt(x, y)

			sumd := 0
			for _, p := range points {
				sumd += manhattanDist(p, check)
			}

			if sumd < than {
				n++
			}
		}
	}
	return n
}

func manhattanDist(p, q Point) int {
	dx := p.X - q.X
	if dx < 0 {
		dx = -dx
	}

	dy := p.Y - q.Y
	if dy < 0 {
		dy = -dy
	}

	return dx + dy
}
