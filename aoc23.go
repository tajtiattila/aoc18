package main

import (
	"fmt"
	"log"
	"sort"

	"github.com/tajtiattila/aoc18/bitset"
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

func findbest23(src []nanobot.Bot) {
	type boti struct {
		c nanobot.MPoint
		r int

		bb nanobot.MBox
	}

	var bots []boti
	for _, bot := range src {
		bots = append(bots, boti{
			c:  nanobot.MPt(bot.X, bot.Y, bot.Z),
			r:  bot.Radius,
			bb: bot.MBox(),
		})
	}

	splits := make([][]int, 4)
	for axis := range splits {
		var v []int
		for _, bot := range bots {
			v = append(v, bot.bb.Min[axis], bot.bb.Max[axis])
		}
		sort.Ints(v)

		// dedup
		j, lastc := 1, v[0]
		for _, c := range v[1:] {
			if c != lastc {
				v[j] = c
				j++
				lastc = c
			}
		}

		v = v[:j]
		splits[axis] = v[:j]
	}

	type axisrange struct {
		lo, hi int // index into splits[axis]
	}
	var bounds [4]axisrange
	for axis := range bounds {
		bounds[axis].lo = 0
		bounds[axis].hi = len(splits[axis]) - 1
	}

	var best struct {
		x, y, z int
		count   int
		dist    int
	}

	yield := func(rng [4]axisrange, count int) {
		if count > best.count {
			best.count = count
			var bb nanobot.MBox
			for axis := 0; axis < 4; axis++ {
				bb.Min[axis] = splits[axis][rng[axis].lo]
				bb.Max[axis] = splits[axis][rng[axis].hi]
			}
			if bb.Empty() {
				panic("empty result")
			}
			x, y, z := bb.Min.Coords()
			best.x, best.y, best.z = x, y, z
			best.dist = x + y + z
		}
	}

	var rec func(axis int, rng [4]axisrange, active bitset.Bitset)

	rec = func(axis int, rng [4]axisrange, active bitset.Bitset) {

		cansplit := false
		for i := 0; i < 4; i++ {
			axis = (axis + 1) % 4
			if rng[axis].lo+1 < rng[axis].hi {
				cansplit = true
				break
			}
		}

		if !cansplit {
			yield(rng, active.Count())
			return
		}

		if active.Count() < best.count {
			return
		}

		lo, hi := rng[axis].lo, rng[axis].hi
		mid := (lo + hi) / 2
		split := splits[axis][mid]

		var nlo, nhi bitset.Bitset
		for i, bot := range bots {
			if active.Get(i) {
				if bot.bb.Min[axis] < split {
					nlo.Set(i)
				}
				if split < bot.bb.Max[axis] {
					nhi.Set(i)
				}
			}
		}

		rlo, rhi := rng, rng
		rlo[axis].hi = mid
		rhi[axis].lo = mid
		if nlo.Count() > nhi.Count() {
			rec(axis, rlo, nlo)
			rec(axis, rhi, nhi)
		} else {
			rec(axis, rhi, nhi)
			rec(axis, rlo, nlo)
		}
	}

	rec(0, bounds, bitset.Ones(len(bots)))

	fmt.Println("23/2:", best.dist)
}
