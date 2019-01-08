package main

import (
	"fmt"
	"log"
)

func teleport23() {
	v := getnanobots()

	fmt.Println("23/1:", maxinrange(v, bestradius(v)))
}

type nanobot struct {
	x, y, z int64
	r       int64
}

func getnanobots() []nanobot {
	lines := PuzzleInputLines(23)

	var v []nanobot
	for i, l := range lines {
		var n nanobot
		//pos=<1,1,1>, r=1
		_, err := fmt.Sscanf(l, "pos=<%d,%d,%d>, r=%d", &n.x, &n.y, &n.z, &n.r)
		if err != nil {
			log.Fatalf("parse line %d: %v", i+1, err)
		}

		v = append(v, n)
	}
	return v
}

func bestradius(v []nanobot) int {
	var rmax int64
	var imax int
	for i, n := range v {
		if n.r > rmax {
			imax, rmax = i, n.r
		}
	}
	return imax
}

func maxinrange(v []nanobot, i int) int {

	n := v[i]

	inrange := 0
	for _, m := range v {
		d := abs64(n.x-m.x) + abs64(n.y-m.y) + abs64(n.z-m.z)
		if d <= n.r {
			inrange++
		}
	}

	return inrange
}

func abs64(x int64) int64 {
	if x >= 0 {
		return x
	}
	return -x
}
