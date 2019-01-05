// reservoir research
package resrsrch

import (
	"bytes"
	"fmt"
	"io"
	"strings"

	"github.com/pkg/errors"
)

type GroundSlice struct {
	// x offset for input coordinates
	// any other x coordinate below is 0-based
	xofs int
	bbox bbox

	dx, dy int // grid dimensions, dx == ystride
	grid   []tile
}

type tile byte

const (
	tilesand tile = iota
	tileclay

	tilewater // static water
	tileflow  // flowing water
)

type bbox struct {
	ix, ax int // x min/max (inclusive)
	iy, ay int // y min/max (inclusive)
}

func (b *bbox) add(o bbox) {
	if o.ix < b.ix {
		b.ix = o.ix
	}
	if o.iy < b.iy {
		b.iy = o.iy
	}
	if o.ax > b.ax {
		b.ax = o.ax
	}
	if o.ay > b.ay {
		b.ay = o.ay
	}
}

func ParseGroundSlice(src []string) (*GroundSlice, error) {
	var cb []bbox
	for i, line := range src {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		var h, v bbox
		_, errh := fmt.Sscanf(line, "y=%d, x=%d..%d", &h.iy, &h.ix, &h.ax)
		_, errv := fmt.Sscanf(line, "x=%d, y=%d..%d", &v.ix, &v.iy, &v.ay)
		if errh != nil && errv != nil {
			return nil, errors.Errorf("cannot scan line %d", i+1)
		}

		if errh == nil {
			h.ay = h.iy
			cb = append(cb, h)
		} else {
			v.ax = v.ix
			cb = append(cb, v)
		}
	}

	if len(cb) == 0 {
		return nil, errors.New("empty groundslice")
	}

	bb := cb[0]
	for _, b := range cb[1:] {
		bb.add(b)
	}

	const hspace = 2
	const vspace = 2

	dy := bb.ay + 1 + vspace
	dx := (bb.ax - bb.ix) + 1 + 2*hspace

	gs := &GroundSlice{
		xofs: hspace - bb.ix,

		dx:   dx,
		dy:   dy,
		grid: make([]tile, dx*dy),

		bbox: bb,
	}

	for _, b := range cb {
		gs.addclaybox(b)
	}

	return gs, nil
}

// ofs calculates offset from x/y coordinates
func (gs *GroundSlice) ofs(x, y int) int { return x + gs.xofs + y*gs.dx }

// yofs calculates y from offset ofs
func (gs *GroundSlice) yofs(ofs int) int {
	if ofs < 0 {
		return -1
	}
	return ofs / gs.dx
}

func (gs *GroundSlice) addclaybox(b bbox) {
	ofs := gs.ofs(b.ix, b.iy)
	for y := b.iy; y <= b.ay; y++ {
		o := ofs
		ofs += gs.dx
		for x := b.ix; x <= b.ax; x++ {
			gs.grid[o] = tileclay
			o++
		}
	}
}

type FloodStat struct {
	Static int // settled water
	Flow   int // flowing water
}

func (fs FloodStat) Total() int {
	return fs.Static + fs.Flow
}

func (gs *GroundSlice) Flood(x, y int, w io.Writer) FloodStat {
	if gs.bbox.ay < y {
		return FloodStat{}
	}

	if x < gs.bbox.ix || gs.bbox.ax < x {
		// todo inf
		return FloodStat{}
	}

	ofs := gs.ofs(x, y)
	if gs.grid[ofs] == tileclay {
		return FloodStat{}
	}

	sim := simstate{
		p: make([]tile, len(gs.grid)),
	}
	copy(sim.p, gs.grid)

	if w != nil {
		gs.dumpSim(w, &sim, "Start")
	}

	iter := 0
	lastwater := -1
	for sim.nwater > lastwater {
		lastwater = sim.nwater
		gs.flowdown(&sim, x, y)
		if w != nil {
			gs.dumpSim(w, &sim, fmt.Sprintf("Iteration #%d", iter+1))
			iter++
		}
	}

	var fs FloodStat

	si := gs.bbox.iy * gs.dx
	ei := (gs.bbox.ay + 1) * gs.dx
	for _, t := range sim.p[si:ei] {
		switch t {
		case tilewater:
			fs.Static++
		case tileflow:
			fs.Flow++
		}
	}

	return fs
}

func (gs *GroundSlice) dumpSim(w io.Writer, sim *simstate, header string) {
	fmt.Fprintln(w, header)
	var buf bytes.Buffer
	for i, t := range sim.p {
		var r byte
		switch t {
		case tilesand:
			r = '.'
		case tileclay:
			r = '#'
		case tilewater:
			r = '~'
		case tileflow:
			r = '|'
		default:
			r = '?'
		}
		buf.WriteByte(r)

		if (i % gs.dx) == gs.dx-1 {
			fmt.Fprintln(w, buf.String())
			buf.Reset()
		}
	}
}

// flowblock reports if tile t blocks flow
func flowblock(t tile) bool {
	return t == tileclay || t == tilewater
}

type simstate struct {
	p []tile

	nwater int
}

func (gs *GroundSlice) flowdown(sim *simstate, x, y int) {
	if gs.bbox.ay < y {
		return
	}

	ofs := gs.ofs(x, y)
	if flowblock(sim.p[ofs]) {
		panic("logic error; must have stopped above")
	}

	sim.p[ofs] = tileflow

	down := ofs + gs.dx
	if !flowblock(sim.p[down]) {
		// try to fill what's below
		gs.flowdown(sim, x, y+1)

		if !flowblock(sim.p[down]) {
			// still not blocked
			return
		}
	}

	// try to fill horizontal slice, or flow out
	lx, lstop := gs.flowhorz(sim, x, y, -1)
	rx, rstop := gs.flowhorz(sim, x, y, +1)
	if !(lstop && rstop) {
		// has outflow
		return
	}

	// fill reservoir slice
	for wx := lx + 1; wx < rx; wx++ {
		sim.p[gs.ofs(wx, y)] = tilewater
		sim.nwater++
	}
}

func (gs *GroundSlice) flowhorz(sim *simstate, x, y, dx int) (rx int, stopped bool) {
	vstride := gs.dx

	x += dx
	for ofs := gs.ofs(x, y); gs.yofs(ofs) == y; x, ofs = x+dx, ofs+dx {
		if flowblock(sim.p[ofs]) {
			// found wall
			return x, true
		}

		sim.p[ofs] = tileflow

		down := ofs + vstride
		if !flowblock(sim.p[down]) {
			// outflow
			gs.flowdown(sim, x, y+1)

			if !flowblock(sim.p[down]) {
				return x, false
			}
		}
	}

	fmt.Println("at", x, y)

	panic("horizontal outflow")
}
