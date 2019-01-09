package constellation

import "fmt"

const Dim = 4

type Point [Dim]int

func Pt(x, y, z, w int) Point {
	return Point{x, y, z, w}
}

func (p Point) Add(q Point) Point {
	return Point{
		p[0] + q[0],
		p[1] + q[1],
		p[2] + q[2],
		p[3] + q[3],
	}
}

func (p Point) Dist(q Point) int {
	return abs(p[0]-q[0]) +
		abs(p[1]-q[1]) +
		abs(p[2]-q[2]) +
		abs(p[3]-q[3])
}

func abs(i int) int {
	if i >= 0 {
		return i
	}
	return -i
}

func ParsePoint(s string) (Point, error) {
	var p Point
	_, err := fmt.Sscanf(s, "%d,%d,%d,%d", &p[0], &p[1], &p[2], &p[3])
	return p, err
}

func ParsePoints(v []string) ([]Point, error) {
	var pts []Point

	for _, s := range v {
		if s == "" {
			continue
		}
		p, err := ParsePoint(s)
		if err != nil {
			return pts, err
		}

		pts = append(pts, p)
	}

	return pts, nil
}

type space struct {
	star []star
}

type star struct {
	p Point

	next []int // indices into space.star

	cidx int // constellation index
}

func Constellations(points []Point, d int) [][]int {

	t := NewTree(points)

	taken := make([]bool, len(points))

	var result [][]int

	addConst := func(i int) {
		var pts []int

		taken[i] = true
		alive := map[int]struct{}{
			i: struct{}{},
		}

		for len(alive) > 0 {
			for i = range alive {
				break
			}
			delete(alive, i)
			pts = append(pts, i)

			t.WalkNeighbors(points[i], d, func(e Elem) {
				if !taken[e.I] {
					alive[e.I] = struct{}{}
					taken[e.I] = true
				}
			})
		}

		result = append(result, pts)
	}

	for i := range points {
		if !taken[i] {
			addConst(i)
		}
	}

	return result
}

func canMerge(a, b []Point, d int) bool {
	for _, pa := range a {
		for _, pb := range a {
			if pa.Dist(pb) <= d {
				return true
			}
		}
	}
	return false
}
