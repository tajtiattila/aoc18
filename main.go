package main

import (
	"flag"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strconv"
)

var verbose bool

var verbosew io.Writer

func main() {
	flag.BoolVar(&verbose, "v", false, "verbose mode")
	flag.Parse()

	if verbose {
		verbosew = os.Stdout
	} else {
		verbosew = ioutil.Discard
	}

	type puzzle struct {
		n int
		f func()
	}

	var puzzles []puzzle
	pm := make(map[int]func())

	add := func(n int, f func()) {
		puzzles = append(puzzles, puzzle{n: n, f: f})
		pm[n] = f
	}

	add(16, wristdevhack)
	add(17, reservoirresearch)
	add(18, collectlumber)
	add(19, wristdev19)
	add(20, facilitymaxdoors)

	if flag.NArg() == 0 {
		for _, p := range puzzles {
			p.f()
		}
	} else {
		for _, arg := range flag.Args() {
			n, _ := strconv.Atoi(arg)
			if f, ok := pm[n]; ok {
				f()
			} else {
				log.Printf("unknown puzzle: %v", arg)
			}
		}
	}
}
