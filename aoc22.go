package main

import (
	"fmt"

	"github.com/tajtiattila/aoc18/modemaze"
)

func modemaze22() {
	const (
		// puzzle input:
		// depth: 11109
		// target: 9,731

		depth  = 11109
		tx, ty = 9, 731
	)

	m := modemaze.New(depth, tx, ty)

	// 22/1: 5836 too low

	//m.Write(os.Stdout, 800, 800)

	fmt.Println("22/1:", m.RiskLevel())
	fmt.Println("22/2:", m.PathDuration())
}
