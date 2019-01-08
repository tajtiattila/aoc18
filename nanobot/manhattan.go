package nanobot

// point in manhattan space
type MPoint [4]int

func MPt(x, y, z int) MPoint {
	return MPoint{
		-x + y + z,
		x - y + z,
		x + y - z,
		x + y + z,
	}
}

func (p MPoint) mcoords() (xm, ym, zm, wm int) {
	return p[0], p[1], p[2], p[3]
}

/*

xm: -x + y + z,
ym: x - y + z,
zm: x + y - z,
wm: x + y + z,

xm+ym= (-x+y+z)+(x-y+z) = 2z

zm+wm= (x+y-z)+(x+y+z) = 2x+2y

*/

func (p MPoint) Coords() (x, y, z int) {
	x = (p[1] + p[2]) / 2
	y = (p[0] + p[2]) / 2
	z = (p[0] + p[1]) / 2
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
	return b.Min[0] >= b.Max[0] ||
		b.Min[1] >= b.Max[1] ||
		b.Min[2] >= b.Max[2] ||
		b.Min[3] >= b.Max[3]
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
			min(a.Min[0], b.Min[0]),
			min(a.Min[1], b.Min[1]),
			min(a.Min[2], b.Min[2]),
			min(a.Min[3], b.Min[3]),
		},
		Max: MPoint{
			max(a.Max[0], b.Max[0]),
			max(a.Max[1], b.Max[1]),
			max(a.Max[2], b.Max[2]),
			max(a.Max[3], b.Max[3]),
		},
	}
}

func (a MBox) Intersect(b MBox) MBox {
	if a.Min[0] < b.Min[0] {
		a.Min[0] = b.Min[0]
	}
	if a.Max[0] > b.Max[0] {
		a.Max[0] = b.Max[0]
	}

	if a.Min[1] < b.Min[1] {
		a.Min[1] = b.Min[1]
	}
	if a.Max[1] > b.Max[1] {
		a.Max[1] = b.Max[1]
	}

	if a.Min[2] < b.Min[2] {
		a.Min[2] = b.Min[2]
	}
	if a.Max[2] > b.Max[2] {
		a.Max[2] = b.Max[2]
	}

	if a.Min[3] < b.Min[3] {
		a.Min[3] = b.Min[3]
	}
	if a.Max[3] > b.Max[3] {
		a.Max[3] = b.Max[3]
	}

	if a.Empty() {
		return MBox{}
	}

	return a
}

const (
	maxuint = ^uint(0)
	maxint  = int(maxuint >> 1)
)

func (b MBox) MinPoint() (x, y, z int, ok bool) {
	if b.Empty() {
		return 0, 0, 0, false
	}

	lo, hi := 0, maxint/4
	for lo+1 != hi {
		mid := (hi + lo) / 2
		t := Equidist(0, 0, 0, mid).Intersect(b)
		if t.Empty() {
			lo = mid
		} else {
			hi = mid
		}
	}

	cross := Equidist(0, 0, 0, hi+1).Intersect(b)

	if cross.Empty() {
		panic("impossible")
	}

	ok = false
	var mx, my, mz int
	cross.WalkPoints(func(x, y, z int) {
		if !ok || x < mx || (x == mx && y < my) || (x == mx && y == my && z < mz) {
			mx, my, mz = x, y, z
			ok = true
		}
	})
	return mx, my, mz, ok
}

func (bb MBox) WalkPoints(f func(x, y, z int)) {
	for xm := bb.Min[0]; xm < bb.Max[0]; xm++ {
		for ym := bb.Min[1]; ym < bb.Max[1]; ym++ {
			for zm := bb.Min[2]; zm < bb.Max[2]; zm++ {
				for wm := bb.Min[3]; wm < bb.Max[3]; wm++ {
					p := MPoint{xm, ym, zm, wm}
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
	for xm := bb.Min[0]; xm < bb.Max[0]; xm++ {
		for ym := bb.Min[1]; ym < bb.Max[1]; ym++ {
			for zm := bb.Min[2]; zm < bb.Max[2]; zm++ {
				for wm := bb.Min[3]; wm < bb.Max[3]; wm++ {
					p := MPoint{xm, ym, zm, wm}
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
