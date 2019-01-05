package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"

	"github.com/tajtiattila/aoc18/lumbercoll"
)

func collectlumber() {
	a, err := lumbercoll.ParseArea(PuzzleInputLines(18))
	if err != nil {
		log.Fatal(err)
	}

	const firstStop = 10
	a.Step(firstStop)

	fmt.Println("18/1:", a.ResourceValue())

	const nline = 10

	var w io.Writer
	if verbose {
		w = os.Stdout
	} else {
		w = ioutil.Discard
	}

	const maxSim = 1000

	var deltas []int
	lastrv := a.ResourceValue()
	for i := firstStop; i < maxSim; i++ {
		if i%nline == 0 {
			fmt.Fprintf(w, "\n%4d  ", i)
		}
		a.Step(1)
		rv := a.ResourceValue()
		delta := rv - lastrv
		deltas = append(deltas, delta)
		fmt.Fprintf(w, "%8d", delta)
		lastrv = rv
	}

	n := findEndRepeat(deltas)
	rpt := deltas[len(deltas)-n:]
	sumrpt := sumIntSlice(rpt)
	fmt.Fprintf(w, "\nend repeat=%d, sum=%d\n", n, sumrpt)

	const wantSim = 1000000000
	const togo = wantSim - maxSim

	rv := int64(lastrv)
	rv += int64(togo) * int64(sumrpt) // sum of repeating values

	rem := int(togo % int64(n)) // remaining indices
	for _, add := range rpt[:rem] {
		rv += int64(add)
	}

	fmt.Println("18/2:", rv)
}

func findEndRepeat(v []int) int {
	e := len(v)
	for n := 1; n < e/2; n++ {
		if equalIntSlice(v[e-2*n:e-n], v[e-n:]) {
			return n
		}
	}
	return 0
}

func equalIntSlice(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func sumIntSlice(v []int) int {
	sum := 0
	for _, i := range v {
		sum += i
	}
	return sum
}
