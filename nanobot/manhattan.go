package nanobot

// point in manhattan space
type MPoint struct {
	Xm, Ym, Zm int
}

func Mpt(x, y, z int) MPoint {
	return MPoint{
		Xm: x + y + z,
		Ym: -x - y + z,
		Zm: -x + y - z,
	}
}

func (p MPoint) Valid() bool {
	return (p.Xm+p.Ym+p.Zm)%4 == 0
}

func (p MPoint) mcoords() (xm, ym, zm int) {
	return p.Xm, p.Ym, p.Zm
}

func (p MPoint) Coords() (x, y, z int) {
	x = -(p.Ym + p.Zm) / 2
	y = (p.Xm + p.Zm) / 2
	z = (p.Xm + p.Ym) / 2
	return
}

// box in manhattan space
type MBox struct {
	Min, Max MPoint
}

func (b MBox) Empty() bool {
	return b.Min.Xm >= b.Max.Xm ||
		b.Min.Ym >= b.Max.Ym ||
		b.Min.Zm >= b.Max.Zm
}

func (b MBox) Extend(o MBox) MBox {
	if b.Empty() {
		return o
	}
	if o.Empty() {
		return b
	}

	return MBox{
		Min: MPoint{
			min(b.Min.Xm, o.Min.Xm),
			min(b.Min.Ym, o.Min.Ym),
			min(b.Min.Zm, o.Min.Zm),
		},
		Max: MPoint{
			max(b.Max.Xm, o.Max.Xm),
			max(b.Max.Ym, o.Max.Ym),
			max(b.Max.Zm, o.Max.Zm),
		},
	}
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

func Equidist(x, y, z, r int) MBox {
	if r <= 0 {
		return MBox{}
	}
	xc, yc, zc := Mpt(x, y, z).mcoords()
	return MBox{
		Min: MPoint{
			xc - r,
			yc - r,
			zc - r,
		},
		Max: MPoint{
			xc + r + 1,
			yc + r + 1,
			zc + r + 1,
		},
	}
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
