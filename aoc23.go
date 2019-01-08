package main

import (
	"fmt"
	"log"

	"github.com/tajtiattila/aoc18/nanobot"
)

func teleport23() {
	v := getnanobots()

	maxr := bestradius(v)
	if verbose {
		b := v[maxr]
		fmt.Printf("maxr: (%d) pos=<%d,%d,%d> r=%d\n", maxr, b.X, b.Y, b.Z, b.Radius)
	}
	fmt.Println("23/1:", maxinrange(v, maxr))
	findbest23(v)
}

func getnanobots() []nanobot.Bot {
	lines := PuzzleInputLines(23)

	var v []nanobot.Bot
	for i, l := range lines {
		var n nanobot.Bot
		//pos=<1,1,1>, r=1
		_, err := fmt.Sscanf(l, "pos=<%d,%d,%d>, r=%d", &n.X, &n.Y, &n.Z, &n.Radius)
		if err != nil {
			log.Fatalf("parse line %d: %v", i+1, err)
		}

		v = append(v, n)
	}
	return v
}

func bestradius(v []nanobot.Bot) int {
	var rmax int
	var imax int
	for i, n := range v {
		if n.Radius > rmax {
			imax, rmax = i, n.Radius
		}
	}
	return imax
}

func maxinrange(v []nanobot.Bot, i int) int {

	n := v[i]

	inrange := 0
	for _, m := range v {
		d := abs(n.X-m.X) + abs(n.Y-m.Y) + abs(n.Z-m.Z)
		if d <= n.Radius {
			inrange++
		}
	}

	return inrange
}

func abs(x int) int {
	if x >= 0 {
		return x
	}
	return -x
}

type findbestinf struct {
	bots []nanobot.Bot

	best []int
}

func findbest23(bots []nanobot.Bot) {
	inf := findbestinf{
		bots: bots,
	}
	for i, b := range bots {
		fmt.Println(i)
		box := nanobot.Equidist(b.X, b.Y, b.Z, b.Radius)
		findbest23x(&inf, box, []int{i})
	}
}

func findbest23x(inf *findbestinf, box nanobot.MBox, v []int) {
	n := len(v)
	starti := v[n-1] + 1

	if rest := len(inf.bots) - starti; len(v)+rest < len(inf.best) {
		return // can't get better
	}

	added := false
	for i := starti; i < len(inf.bots); i++ {
		b := inf.bots[i]
		c := nanobot.Equidist(b.X, b.Y, b.Z, b.Radius)
		cross := box.Intersect(c)
		if !cross.Empty() {
			v = append(v[:n], i)
			findbest23x(inf, cross, v)
			added = true
		}
	}
	v = v[:n]
	if !added && len(v) > len(inf.best) {
		inf.best = make([]int, len(v))
		copy(inf.best, v)
		fmt.Println(" ", len(v))
		fmt.Println(" ", v)
		box.WalkPoints(func(x, y, z int) {
			fmt.Println(" ", x, y, z)
		})
	}
}
