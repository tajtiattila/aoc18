package main

type CutSpec struct {
	ID int // claimant

	// cut spec
	X  int
	Y  int
	Dx int
	Dy int
}

func FindCutSpecOverlap(specs []CutSpec, minovl int) int {
	dx, dy := maxDim(specs)
	f := newFabric(dx, dy)
	for _, cs := range specs {
		f.addCutSpec(cs)
	}

	res := 0
	for _, v := range f.data {
		if v >= minovl {
			res++
		}
	}
	return res
}

func maxDim(v []CutSpec) (dx, dy int) {
	for _, cs := range v {
		cx, cy := cs.X+cs.Dx, cs.Y+cs.Dy
		if cx > dx {
			dx = cx
		}
		if cy > dy {
			dy = cy
		}
	}
	return dx, dy
}

type fabric struct {
	dx, dy int
	data   []int
}

func newFabric(dx, dy int) *fabric {
	return &fabric{
		dx:   dx,
		dy:   dy,
		data: make([]int, dx*dy),
	}
}

func (f *fabric) ofs(x, y int) int { return x + y*f.dx }

func (f *fabric) add(x, y int, v int) {
	f.data[f.ofs(x, y)] += v
}

func (f *fabric) addCutSpec(cs CutSpec) {
	for x := 0; x < cs.Dx; x++ {
		for y := 0; y < cs.Dy; y++ {
			f.add(cs.X+x, cs.Y+y, 1)
		}
	}
}

func (f *fabric) get(x, y int) int {
	return f.data[f.ofs(x, y)]
}
