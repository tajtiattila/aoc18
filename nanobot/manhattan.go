package nanobot

// point in manhattan space
type MPoint struct {
	Xm, Ym, Zm, Wm int
}

func MPt(x, y, z int) MPoint {
	return MPoint{
		Xm: -x + y + z,
		Ym: x - y + z,
		Zm: x + y - z,
		Wm: x + y + z,
	}
}

func (p MPoint) mcoords() (xm, ym, zm, wm int) {
	return p.Xm, p.Ym, p.Zm, p.Wm
}

/*

Xm: -x + y + z,
Ym: x - y + z,
Zm: x + y - z,
Wm: x + y + z,

Xm+Ym= (-x+y+z)+(x-y+z) = 2z

Zm+Wm= (x+y-z)+(x+y+z) = 2x+2y

*/

func (p MPoint) Coords() (x, y, z int) {
	x = (p.Ym + p.Zm) / 2
	y = (p.Xm + p.Zm) / 2
	z = (p.Xm + p.Ym) / 2
	return
}

func (p MPoint) Valid() bool {
	x, y, z := p.Coords()
	return MPt(x, y, z) == p
}

// box in manhattan space
type MBox struct {
	Min, Max MPoint
}

func Equidist(x, y, z, r int) MBox {
	if r <= 0 {
		return MBox{}
	}
	xc, yc, zc, wc := MPt(x, y, z).mcoords()
	return MBox{
		Min: MPoint{
			xc - r,
			yc - r,
			zc - r,
			wc - r,
		},
		Max: MPoint{
			xc + r + 1,
			yc + r + 1,
			zc + r + 1,
			wc + r + 1,
		},
	}
}

func (b MBox) Volume() int {
	xi, yi, zi, wi := b.Min.mcoords()
	xa, ya, za, wa := b.Max.mcoords()

	dx, dy, dz, dw := xa-xi, ya-yi, za-zi, wa-wi
	return dx * dy * dz * dw
}

func (b MBox) Empty() bool {
	return b.Min.Xm >= b.Max.Xm ||
		b.Min.Ym >= b.Max.Ym ||
		b.Min.Zm >= b.Max.Zm ||
		b.Min.Wm >= b.Max.Wm
}

func (a MBox) Extend(b MBox) MBox {
	if a.Empty() {
		return b
	}
	if b.Empty() {
		return a
	}

	return MBox{
		Min: MPoint{
			min(a.Min.Xm, b.Min.Xm),
			min(a.Min.Ym, b.Min.Ym),
			min(a.Min.Zm, b.Min.Zm),
			min(a.Min.Wm, b.Min.Wm),
		},
		Max: MPoint{
			max(a.Max.Xm, b.Max.Xm),
			max(a.Max.Ym, b.Max.Ym),
			max(a.Max.Zm, b.Max.Zm),
			max(a.Max.Wm, b.Max.Wm),
		},
	}
}

func (a MBox) Intersect(b MBox) MBox {
	if a.Min.Xm < b.Min.Xm {
		a.Min.Xm = b.Min.Xm
	}
	if a.Max.Xm > b.Max.Xm {
		a.Max.Xm = b.Max.Xm
	}

	if a.Min.Ym < b.Min.Ym {
		a.Min.Ym = b.Min.Ym
	}
	if a.Max.Ym > b.Max.Ym {
		a.Max.Ym = b.Max.Ym
	}

	if a.Min.Zm < b.Min.Zm {
		a.Min.Zm = b.Min.Zm
	}
	if a.Max.Zm > b.Max.Zm {
		a.Max.Zm = b.Max.Zm
	}

	if a.Min.Wm < b.Min.Wm {
		a.Min.Wm = b.Min.Wm
	}
	if a.Max.Wm > b.Max.Wm {
		a.Max.Wm = b.Max.Wm
	}

	if a.Empty() {
		return MBox{}
	}

	return a
}

func (b MBox) MinPoint(upto int) (x, y, z int, ok bool) {
	var cross MBox
	for i := 1; i < upto; i++ {
		c := Equidist(0, 0, 0, i).Intersect(b)
		if !c.Empty() {
			cross = c
			break
		}
	}

	if cross.Empty() {
		return 0, 0, 0, false
	}

	first := true
	var mx, my, mz int
	cross.WalkPoints(func(x, y, z int) {
		if first || x < mx || (x == mx && y < my) || (x == mx && y == my && z < mz) {
			mx, my, mz = x, y, z
			first = false
		}
	})
	return mx, my, mz, true
}

func (bb MBox) WalkPoints(f func(x, y, z int)) {
	for xm := bb.Min.Xm; xm < bb.Max.Xm; xm++ {
		for ym := bb.Min.Ym; ym < bb.Max.Ym; ym++ {
			for zm := bb.Min.Zm; zm < bb.Max.Zm; zm++ {
				for wm := bb.Min.Wm; wm < bb.Max.Wm; wm++ {
					p := MPoint{Xm: xm, Ym: ym, Zm: zm, Wm: wm}
					x, y, z := p.Coords()
					if MPt(x, y, z) == p {
						f(x, y, z)
					}
				}
			}
		}
	}
}

func (bb MBox) NumPoints() int {
	n := 0
	for xm := bb.Min.Xm; xm < bb.Max.Xm; xm++ {
		for ym := bb.Min.Ym; ym < bb.Max.Ym; ym++ {
			for zm := bb.Min.Zm; zm < bb.Max.Zm; zm++ {
				for wm := bb.Min.Wm; wm < bb.Max.Wm; wm++ {
					p := MPoint{Xm: xm, Ym: ym, Zm: zm, Wm: wm}
					if p.Valid() {
						n++
					}
				}
			}
		}
	}
	return n
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

/*

2D manhattan space:
 ym | xm 0  1  2  3  4
 -2           0,2
 -1        0,1   1,2
  0     0,0   1,1   2,2
  1        1,0   2,1
  2           2,0

Point (1,1) in manhattan space is (2,0).
With r=1 its points in manhattan space are in
([xm-r..xm+r], [ym-r..xm+r]) = (1..3, -1..1)


3D manhattan space:

xm = (x+y+z)
ym = (-x+y-z)
zm = (-x-y+z)

xm+ym = (x+y+z)+(-x-y+z) = 2z
xm+zm = (x+y+z)+(-x+y-z) = 2y
ym+zm = (-x-y+z)+(-x+y-z) = -2x

xm+ym+zm = (-x+y+z)

Only coordinates where (xm+ym+zm)%2 == 0 are valid

*/

/* scrap

xc=(2,2,2), r=1 points are

cart.   manh.
2,2,2   6,-2,-2
1,2,2   5,-1,-1
3,2,2   7,-3,-3
2,1,2   5,-3,-1
2,3,2   7,-1,-3
2,2,1   5,-1,-3
2,2,3   7,-3,-1

!!!!


3D manhattan space:

xm = (x-y-z)
ym = (-x+y-z)
zm = (-x-y+z)

xm+ym = (x-y-z)+(
Only coordinates where (xm+ym+zm)%2 == 0 are valid


xc=(2,2,2), r=1 points are

cart.   manh.
2,2,2   -2,-2,-2
1,2,2   -3,-1,-1
3,2,2   -1,-3,-3
2,1,2   -1,-3,-1
2,3,2   -3,-1,-3
2,2,1   -1,-1,-3
2,2,3   -3,-3,-1



Bounding box of all points from (1,1) with manhattan radius 1

*/
