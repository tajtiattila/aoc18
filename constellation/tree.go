package constellation

import (
	"sort"
)

var ZeroBox Box

type Box struct {
	Min, Max Point
}

func (b *Box) Add(p Point) {
	for i := range p {
		if p[i] < b.Min[i] {
			b.Min[i] = p[i]
		}
		if b.Max[i] <= p[i] {
			b.Max[i] = p[i] + 1
		}
	}
}

func (b Box) In(p Point) bool {
	for i := range p {
		if !(b.Min[i] <= p[i] || p[i] < b.Max[i]) {
			return false
		}
	}
	return true
}

func (b Box) Intersect(c Box) Box {
	for i := range b.Min {
		if b.Min[i] < c.Min[i] {
			b.Min[i] = c.Min[i]
		}
		if c.Max[i] < b.Max[i] {
			b.Max[i] = c.Max[i]
		}
		if b.Max[i] <= b.Min[i] {
			return ZeroBox
		}
	}
	return b
}

type Tree struct {
	Box

	Elem []Elem

	Child []Tree
}

type Elem struct {
	P Point
	I int // index
}

func NewTree(points []Point) *Tree {
	if len(points) == 0 {
		return &Tree{}
	}

	v := make([]Elem, len(points))
	for i, p := range points {
		v[i] = Elem{P: p, I: i}
	}

	t := newTree(v)

	const maxelem = 8

	t.splitOnCount(maxelem)

	return t
}

func newTree(elem []Elem) *Tree {
	bb := Box{
		Min: elem[0].P,
		Max: elem[0].P,
	}
	for _, e := range elem {
		bb.Add(e.P)
	}

	return &Tree{
		Box:  bb,
		Elem: elem,
	}
}

func (t *Tree) splitOnCount(n int) {
	if len(t.Elem) < n {
		return
	}
	if !t.split() {
		return
	}

	for i := range t.Child {
		t.Child[i].splitOnCount(n)
	}
}

func (t *Tree) split() bool {
	if len(t.Elem) == 0 {
		return false
	}

	if len(t.Child) != 0 {
		panic("already split")
	}

	// midpoint of split
	var mid Point

	// find optimal split point
	nchild := 1
	cansplit := false
	for axis := range mid {
		var v []int

		for _, e := range t.Elem {
			v = append(v, e.P[axis])
		}
		sort.Ints(v)

		// unique
		lastc := v[0]
		j := 1
		for _, c := range v[1:] {
			if c != lastc {
				v[j], j = c, j+1
				lastc = c
			}
		}
		v = v[:j]

		if len(v) > 1 {
			cansplit = true
		}

		midc := v[len(v)/2]
		mid[axis] = midc
		nchild *= 2
	}

	if !cansplit {
		return false
	}

	vchild := make([][]Elem, nchild)
	for _, e := range t.Elem {
		index, m := 0, 1
		for axis := range mid {
			if mid[axis] <= e.P[axis] {
				index += m
			}
			m *= 2
		}
		vchild[index] = append(vchild[index], e)
	}

	for _, child := range vchild {
		if len(child) != 0 {
			node := newTree(child)
			t.Child = append(t.Child, *node)
		}
	}

	t.Elem = nil

	return true
}

// WalkNeighbors walks neightbors of p
// within manhattan distance d.
func (t *Tree) WalkNeighbors(p Point, d int, f func(e Elem)) {
	var bb Box
	for axis := range p {
		bb.Min[axis] = p[axis] - d
		bb.Max[axis] = p[axis] + d + 1
	}

	t.walkBox(bb, func(e Elem) {
		if e.P.Dist(p) <= d {
			f(e)
		}
	})
}

func (t *Tree) walkBox(box Box, f func(e Elem)) {
	if t.Intersect(box) == ZeroBox {
		return
	}

	if t.Child != nil {
		for i := range t.Child {
			t.Child[i].walkBox(box, f)
		}
		return
	}

	for _, e := range t.Elem {
		if box.In(e.P) {
			f(e)
		}
	}
}
