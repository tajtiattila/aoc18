package main

import (
	"bytes"
	"fmt"
	"sort"

	"github.com/pkg/errors"
)

type mctile byte

const (
	mtEmpty mctile = iota

	// straights
	mtEW // "horizontal"
	mtNS // "vertical"

	// turns
	mtLeft
	mtRight

	mtCross
)

var mctileRuneMap = map[mctile]rune{
	mtEmpty: ' ',
	mtEW:    '-',
	mtNS:    '|',
	mtLeft:  '\\',
	mtRight: '/',
	mtCross: '+',
}

type mcnextkey struct {
	tile mctile
	dir  int
}

var mcnextDirMap map[mcnextkey]int

func init() {
	mcnextDirMap = make(map[mcnextkey]int)
	e := func(tile mctile, d1, nd1, d2, nd2 int) {
		mcnextDirMap[mcnextkey{tile: tile, dir: d1}] = nd1
		mcnextDirMap[mcnextkey{tile: tile, dir: d2}] = nd2
	}

	e(mtEW, 1, 1, 3, 3)
	e(mtNS, 0, 0, 2, 2)

	e(mtLeft, 0, 3, 2, 1)
	e(mtLeft, 1, 2, 3, 0)
	e(mtRight, 0, 1, 2, 3)
	e(mtRight, 1, 0, 3, 2)
}

type MinecartMap struct {
	dx, dy int
	t      []mctile

	cart []minecart

	crash []Point
}

func (m *MinecartMap) ofs(x, y int) int       { return x + y*m.dx }
func (m *MinecartMap) tile(x, y int) mctile   { return m.t[m.ofs(x, y)] }
func (m *MinecartMap) set(x, y int, t mctile) { m.t[m.ofs(x, y)] = t }

type minecart struct {
	x, y    int
	dir     int // 0: north, 1: east, 2: south...
	xop     int // next op at crossing (0: left, 1: straight, 2: right)
	crashed bool
}

func ParseMinecartMap(src []string) (*MinecartMap, error) {
	if len(src) == 0 {
		return nil, errors.New("empty map")
	}

	m := &MinecartMap{
		dx: len(src[0]),
		dy: len(src),
	}

	m.t = make([]mctile, m.dx*m.dy)

	addCart := func(x, y, dir int) {
		m.cart = append(m.cart, minecart{
			x:   x,
			y:   y,
			dir: dir,
		})
	}

	for y, s := range src {
		if len(s) != m.dx {
			return nil, errors.Errorf("invalid line length in line %v", y)
		}

		ofs := y * m.dx
		for x := 0; x < m.dx; x, ofs = x+1, ofs+1 {
			c := s[x]
			switch c {
			case '-', '<', '>':
				m.t[ofs] = mtEW
			case '|', '^', 'v':
				m.t[ofs] = mtNS
			case '+':
				m.t[ofs] = mtCross
			case '/':
				m.t[ofs] = mtRight
			case '\\':
				m.t[ofs] = mtLeft
			case ' ':
				// pass
			default:
				return nil, errors.Errorf("invalid byte %v in line %v", c, y)
			}

			switch c {
			case '^':
				addCart(x, y, 0)
			case '>':
				addCart(x, y, 1)
			case 'v':
				addCart(x, y, 2)
			case '<':
				addCart(x, y, 3)
			}
		}
	}

	return m, nil
}

func (m *MinecartMap) sortcarts() {
	sort.Slice(m.cart, func(i, j int) bool {
		ci := m.cart[i]
		cj := m.cart[j]
		if ci.y != cj.y {
			return ci.y < cj.y
		}
		return ci.x < cj.x
	})
}

func (m *MinecartMap) Tick() {
	m.sortcarts()

	ncrash := len(m.crash)
	for i := range m.cart {
		m.stepcart(i)
	}

	if ncrash == len(m.crash) {
		return
	}

	// remove crashed carts
	j := 0
	for _, c := range m.cart {
		if !c.crashed {
			m.cart[j] = c
			j++
		}
	}
	m.cart = m.cart[:j]
}

func (m *MinecartMap) stepcart(index int) {
	cart := &m.cart[index]
	if cart.crashed {
		return
	}

	switch cart.dir {
	case 0:
		cart.y--
	case 1:
		cart.x++
	case 2:
		cart.y++
	case 3:
		cart.x--
	default:
		panic("lost")
	}

	tile := m.tile(cart.x, cart.y)

	if tile == mtCross {
		cart.dir = (cart.dir + cart.xop + 4 - 1) % 4
		cart.xop = (cart.xop + 1) % 3
	} else {
		nd, ok := mcnextDirMap[mcnextkey{tile: tile, dir: cart.dir}]
		if !ok {
			panic("derailed")
		}
		cart.dir = nd
	}

	for i := range m.cart {
		if i == index {
			continue
		}
		other := &m.cart[i]
		if cart.x == other.x && cart.y == other.y {
			cart.crashed = true
			other.crashed = true
			m.crash = append(m.crash, Pt(cart.x, cart.y))
		}
	}
}

func (m *MinecartMap) Print() {
	buf := &bytes.Buffer{}
	for y := 0; y < m.dy; y++ {
		buf.Reset()
		for x := 0; x < m.dx; x++ {
			buf.WriteRune(m.runeAt(x, y))
		}
		fmt.Println(buf.String())
	}
}

func (m *MinecartMap) runeAt(x, y int) rune {
	for _, c := range m.cart {
		if c.x == x && c.y == y {
			return minecartRune(c.dir)
		}
	}
	return mctileRuneMap[m.tile(x, y)]
}

func minecartRune(dir int) rune {
	switch dir {
	case 0:
		return '^'
	case 1:
		return '>'
	case 2:
		return 'v'
	case 3:
		return '<'
	default:
		return '?'
	}
}
