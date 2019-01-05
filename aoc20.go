package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/tajtiattila/aoc18/gridregexp"
)

func facilitymaxdoors() {
	src := strings.Join(PuzzleInputLines(20), "")

	gr, err := gridregexp.Parse(src)
	if err != nil {
		log.Fatal(err)
	}

	w := verbosew

	m := gr.Map()
	fmt.Fprintln(w, "extent:", m.Bounds())

	fmt.Println("20/1:", m.MaxDoors())
	fmt.Println("20/2:", m.FarRooms(1000))
}
