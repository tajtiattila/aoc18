package main

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/tajtiattila/aoc18/resrsrch"
)

func reservoirresearch() {
	var w io.Writer

	if verbose {
		w = os.Stdout
	}

	gs, err := resrsrch.ParseGroundSlice(PuzzleInputLines(17))
	if err != nil {
		log.Fatal(err)
	}

	stat := gs.Flood(500, 0, w)

	fmt.Println("17/1:", stat.Total())
	fmt.Println("17/2:", stat.Static)
}
