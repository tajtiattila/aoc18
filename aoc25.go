package main

import (
	"fmt"
	"log"

	"github.com/tajtiattila/aoc18/constellation"
)

func constellations() {
	points, err := constellation.ParsePoints(PuzzleInputLines(25))
	if err != nil {
		log.Fatal(err)
	}

	c := constellation.Constellations(points, 3)
	fmt.Println("25/1:", len(c))
}
