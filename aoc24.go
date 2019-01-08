package main

import (
	"fmt"
	"log"
	"sort"

	"github.com/tajtiattila/aoc18/immunesys"
)

func immunesysbattle() {
	pi := OpenPuzzleInput(24)
	defer pi.Close()

	battle, err := immunesys.ParseBattle(pi)
	if err != nil {
		log.Fatal("parse input:", err)
	}

	b := battle.Clone()
	b.Run()

	fmt.Println("24/1:", b.TotalUnitCount())

	const wantWinner = "Immune System"

	lo, hi := 0, 1000
	for {
		b := battle.Clone()
		b.Boost(wantWinner, hi)
		winner, ok := b.Run()
		if verbose {
			fmt.Println(hi, winner, ok)
		}
		if !ok || winner != wantWinner {
			lo, hi = hi, hi*2
		} else {
			break
		}
	}

	needBoost := lo + sort.Search(hi-lo, func(i int) bool {
		boost := lo + i

		b := battle.Clone()
		b.Boost(wantWinner, boost)
		//b.ShowHeader(verbosew)
		winner, ok := b.Run()
		if verbose {
			fmt.Println(boost, winner, ok)
		}
		return ok && winner == wantWinner
	})

	b = battle.Clone()
	b.Boost(wantWinner, needBoost)
	b.Run()

	fmt.Println("24/2:", b.TotalUnitCount())
}
