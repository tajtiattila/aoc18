package lumbercoll

import (
	"io"

	"github.com/pkg/errors"
)

type Area struct {
	dx, dy int // true dimensions

	start   int // start offset of 0,0 coordinate
	vstride int // vertical stride

	m []tile

	last []tile
}

type tile byte

const (
	tileopen = iota
	tiletree
	tileyard
)

const border = 1

func ParseArea(src []string) (*Area, error) {
	if len(src) == 0 {
		return nil, errors.New("empty area")
	}

	dx, dy := len(src[0]), len(src)

	vstride := dx + 2*border

	n := (dy + 2*border) * vstride

	a := &Area{
		dx: dx,
		dy: dy,

		start:   border * (vstride + 1),
		vstride: vstride,

		m: make([]tile, n),

		last: make([]tile, n),
	}

	for y, line := range src {
		if len(line) != dx {
			return nil, errors.Errorf("line %d: invalid length", y+1)
		}

		for x := 0; x < dx; x++ {
			var t tile
			switch line[x] {
			case '.':
				t = tileopen
			case '|':
				t = tiletree
			case '#':
				t = tileyard
			default:
				return nil, errors.Errorf("line %d: invalid rune", y+1)
			}
			a.m[a.ofs(x, y)] = t
		}
	}

	return a, nil
}

func (a *Area) ofs(x, y int) int { return a.start + x + y*a.vstride }

func (a *Area) nbors(ofs int, res []int) []int {
	v := a.vstride
	return append(res,
		ofs-v-1, ofs-v, ofs-v+1,
		ofs-1, ofs+1,
		ofs+v-1, ofs+v, ofs+v+1)
}

func (a *Area) Step(n int) {
	for i := 0; i < n; i++ {
		a.stepone()
	}
}

func (a *Area) ResourceValue() int {
	return a.TreeCount() * a.LumberyardCount()
}

func (a *Area) TreeCount() int       { return a.count(tiletree) }
func (a *Area) LumberyardCount() int { return a.count(tileyard) }

func (a *Area) count(t tile) int {
	n := 0
	line := a.start
	for y := 0; y < a.dy; y++ {
		o := line
		line += a.vstride
		for x := 0; x < a.dx; x, o = x+1, o+1 {
			if a.m[o] == t {
				n++
			}
		}
	}
	return n
}

func (a *Area) Dump(w io.Writer) error {
	buf := make([]byte, a.dx+1)
	buf[a.dx] = '\n'
	line := a.start
	for y := 0; y < a.dy; y++ {
		o := line
		line += a.vstride
		for x := 0; x < a.dx; x, o = x+1, o+1 {
			var c byte
			switch a.m[o] {
			case tileopen:
				c = '.'
			case tiletree:
				c = '|'
			case tileyard:
				c = '#'
			default:
				c = '?'
			}
			buf[x] = c
		}
		_, err := w.Write(buf)
		if err != nil {
			return err
		}
	}
	return nil
}

func (a *Area) stepone() {
	copy(a.last, a.m)

	var buf [8]int

	line := a.start
	for y := 0; y < a.dy; y++ {
		o := line
		line += a.vstride
		for x := 0; x < a.dx; x, o = x+1, o+1 {
			nbors := a.nbors(o, buf[:0])
			switch a.m[o] {

			case tileopen:
				ntrees := 0
				for _, no := range nbors {
					if a.last[no] == tiletree {
						ntrees++
					}
				}
				if ntrees >= 3 {
					a.m[o] = tiletree
				}

			case tiletree:
				nyards := 0
				for _, no := range nbors {
					if a.last[no] == tileyard {
						nyards++
					}
				}
				if nyards >= 3 {
					a.m[o] = tileyard
				}

			case tileyard:
				ntrees, nyards := 0, 0
				for _, no := range nbors {
					switch a.last[no] {
					case tiletree:
						ntrees++
					case tileyard:
						nyards++
					}
				}
				if ntrees == 0 || nyards == 0 {
					a.m[o] = tileopen
				}
			}
		}
	}
}
