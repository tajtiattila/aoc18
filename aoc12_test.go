package main

import "testing"

func TestAoc12(t *testing.T) {
	sample := `initial state: #..#.#..##......###...###

...## => #
..#.. => #
.#... => #
.#.#. => #
.#.## => #
.##.. => #
.#### => #
#.#.# => #
#.### => #
##.#. => #
##.## => #
###.. => #
###.# => #
####. => #`

	pots, err := NewCavePots(sample)
	if err != nil {
		t.Fatal(err)
	}

	simPotN(t, pots, 20)

	got1 := pots.PlantSum()
	want1 := 325
	if got1 != want1 {
		t.Errorf("1: plant sum is %d; want %d", got1, want1)
	}

	pots, err = NewCavePots(input12)
	if err != nil {
		t.Fatal("parse puzzle input:", err)
	}

	simPotN(t, pots, 20)

	t.Logf("1st plant sum: %d", pots.PlantSum())

	pots, err = NewCavePots(input12)
	if err != nil {
		t.Fatal("parse puzzle input:", err)
	}

	const simupto = 1500
	const wantsim = 50000000000

	var delta []int
	var sums []int
	lastsum := pots.PlantSum()
	for i := 0; i < simupto; i++ {
		pots.Step()
		sum := pots.PlantSum()
		d := sum - lastsum
		delta = append(delta, d)
		sums = append(sums, sum)
		lastsum = sum
	}

	lastdelta := delta[len(delta)-1]

	for i := len(delta) - 100; i < len(delta); i++ {
		if delta[i] != lastdelta {
			t.Fatal("delta not repeating")
		}
	}

	final := int64(lastsum) + int64(wantsim-simupto)*int64(lastdelta)

	t.Logf("2nd plant sum: %d", final)
}

func TestPlangAddBorder(t *testing.T) {
	pots, err := NewCavePots(input12)
	if err != nil {
		t.Fatal("parse puzzle input:", err)
	}

	for i := 0; i < 20; i++ {
		pots.Step()
	}

	s1 := pots.PlantSum()
	for grow := 1; grow < 1000; grow *= 2 {
		clone := pots.clone()
		clone.addBorder(grow)
		s2 := clone.PlantSum()
		if s1 != s2 {
			t.Fatalf("planstum changed after grow by %d", grow)
		}
	}
}

func simPotN(t *testing.T, cp *CavePots, n int) {
	for i := 0; i < n; i++ {
		cp.Step()
		//t.Log(cp.Fmt(-3, 39))
	}
}
