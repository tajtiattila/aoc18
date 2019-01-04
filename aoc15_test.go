package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestAoC15PickGoal(t *testing.T) {
	tests := []struct {
		layout string
		sx, sy int // unit
		ex, ey int // chosen goal
	}{
		{`
#######
#E..G.#
#...#.#
#.G.#G#
#######`,
			1, 1, 3, 1,
		},
	}

	for _, tt := range tests {
		gf, err := ParseGoblinFight(tt.layout)
		if err != nil {
			t.Fatal(err)
		}
		got, _ := gf.pickGoal(Pt(tt.sx, tt.sy))
		buf := &bytes.Buffer{}
		gf.DumpFloodMap(buf)
		t.Logf("map:\n%s", buf.String())
		want := Pt(tt.ex, tt.ey)
		if got != want {
			t.Errorf("got %v; want %v", got, want)
		}
	}
}

func TestAoC15Step(t *testing.T) {
	tests := []struct {
		layout           string
		unit, goal, want Point
	}{
		{`
#######
#.E...#
#.....#
#...G.#
#######`,
			Pt(2, 1), Pt(3, 4), Pt(3, 1),
		},
	}

	for _, tt := range tests {
		gf, err := ParseGoblinFight(tt.layout)
		if err != nil {
			t.Fatal(err)
		}
		got := gf.step(tt.unit, tt.goal)
		buf := &bytes.Buffer{}
		gf.DumpFloodMap(buf)
		t.Logf("map:\n%s", buf.String())
		if got != tt.want {
			t.Errorf("got %v; want %v", got, tt.want)
		}
	}
}

func TestAoC15Sim(t *testing.T) {

	type tspec struct {
		after  int
		layout string
	}

	type mtest struct {
		label    string
		outcome  int // 0 means don't check
		withStat bool
		start    string

		elfAttackStrength int

		t []tspec
	}

	sbattle := func(label string, rounds, finalhp, str int, src string) mtest {

		/*
			#######       #######
			#E..EG#       #.E.E.#   E(164), E(197)
			#.#G.E#       #.#E..#   E(200)
			#E.##E#  -->  #E.##.#   E(98)
			#G..#.#       #.E.#.#   E(200)
			#..E#.#       #...#.#
			#######       #######
		*/

		src = strings.TrimSpace(src)
		dx := 0
		for src[dx] == '#' {
			dx++
		}

		var vstart, vfinal []string
		for _, line := range strings.Split(src, "\n") {
			const arrow = 7
			sl, fl := line[:dx], line[dx+arrow:]
			vstart = append(vstart, sl)
			vfinal = append(vfinal, fl)
		}
		start := strings.Join(vstart, "\n")
		final := strings.Join(vfinal, "\n")

		return mtest{
			label:             label,
			outcome:           finalhp * rounds,
			withStat:          true,
			elfAttackStrength: str,
			start:             start,
			t: []tspec{{
				after:  rounds + 1, // +1 for partial round
				layout: final,
			}},
		}
	}

	battle := func(label string, rounds, finalhp int, src string) mtest {
		return sbattle(label, rounds, finalhp, 0, src)
	}

	tests := []mtest{
		mtest{
			label: "aoc1",
			start: `
#########
#G..G..G#
#.......#
#.......#
#G..E..G#
#.......#
#.......#
#G..G..G#
#########`,
			withStat: false,
			t: []tspec{
				{1, `
#########
#.G...G.#
#...G...#
#...E..G#
#.G.....#
#.......#
#G..G..G#
#.......#
#########`,
				},
				{2, `
#########
#..G.G..#
#...G...#
#.G.E.G.#
#.......#
#G..G..G#
#.......#
#.......#
#########`,
				},
				{3, `
#########
#.......#
#..GGG..#
#..GEG..#
#G..G...#
#......G#
#.......#
#.......#
#########`,
				},
			},
		},
		mtest{
			label:    "aoc2",
			outcome:  27730, // 590*47
			withStat: true,
			start: `
#######
#.G...#
#...EG#
#.#.#G#
#..G#E#
#.....#
#######`,
			t: []tspec{
				{1, `
#######   
#..G..#   G(200)
#...EG#   E(197), G(197)
#.#G#G#   G(200), G(197)
#...#E#   E(197)
#.....#   
#######   `,
				},
				{2, `
#######   
#...G.#   G(200)
#..GEG#   G(200), E(188), G(194)
#.#.#G#   G(194)
#...#E#   E(194)
#.....#   
#######   `,
				},
				{23, `
#######   
#...G.#   G(200)
#..G.G#   G(200), G(131)
#.#.#G#   G(131)
#...#E#   E(131)
#.....#   
#######   `,
				},
				{24, `
#######   
#..G..#   G(200)
#...G.#   G(131)
#.#G#G#   G(200), G(128)
#...#E#   E(128)
#.....#   
#######   `,
				},
				{25, `
#######   
#.G...#   G(200)
#..G..#   G(131)
#.#.#G#   G(125)
#..G#E#   G(200), E(125)
#.....#   
#######   `,
				},
				{26, `
#######   
#G....#   G(200)
#.G...#   G(131)
#.#.#G#   G(122)
#...#E#   E(122)
#..G..#   G(200)
#######   `,
				},
				{27, `
#######   
#G....#   G(200)
#.G...#   G(131)
#.#.#G#   G(119)
#...#E#   E(119)
#...G.#   G(200)
#######   `,
				},
				{28, `
#######   
#G....#   G(200)
#.G...#   G(131)
#.#.#G#   G(116)
#...#E#   E(113)
#....G#   G(200)
#######   `,
				},
				{47, `
#######   
#G....#   G(200)
#.G...#   G(131)
#.#.#G#   G(59)
#...#.#   
#....G#   G(200)
#######   `,
				},
			},
		},

		battle("battle1", 37, 982, `
#######       #######
#G..#E#       #...#E#   E(200)
#E#E.E#       #E#...#   E(197)
#G.##.#  -->  #.E##.#   E(185)
#...#E#       #E..#E#   E(200), E(200)
#...E.#       #.....#
#######       #######`),

		battle("battle2", 46, 859, `
#######       #######   
#E..EG#       #.E.E.#   E(164), E(197)
#.#G.E#       #.#E..#   E(200)
#E.##E#  -->  #E.##.#   E(98)
#G..#.#       #.E.#.#   E(200)
#..E#.#       #...#.#   
#######       #######`),

		battle("battle3", 35, 793, `
#######       #######   
#E.G#.#       #G.G#.#   G(200), G(98)
#.#G..#       #.#G..#   G(200)
#G.#.G#  -->  #..#..#
#G..#.#       #...#G#   G(95)
#...E.#       #...G.#   G(200)
#######       #######`),

		battle("battle4", 54, 536, `
#######       #######
#.E...#       #.....#
#.#..G#       #.#G..#   G(200)
#.###.#  -->  #.###.#   
#E#G#G#       #.#.#.#   
#...#G#       #G.G#G#   G(98), G(38), G(200)
#######       #######`),

		battle("battle5", 20, 937, `
#########       #########   
#G......#       #.G.....#   G(137)
#.E.#...#       #G.G#...#   G(200), G(200)
#..##..G#       #.G##...#   G(200)
#...##..#  -->  #...##..#   
#...#...#       #.G.#...#   G(200)
#.G...G.#       #.......#   
#.....G.#       #.......#   
#########       #########`),

		mtest{
			label: "reddit9",
			start: `
########
#.E....#
#......#
#....G.#
#...G..#
#G.....#
########`,
			outcome: 12744, // 531x24
		},
		mtest{
			label: "reddit10",
			start: `
#################
##..............#
##........G.....#
####.....G....###
#....##......####
#...............#
##........GG....#
##.........E..#.#
#####.###...#####
#################`,
			outcome: 14740, // 737x20
		},

		sbattle("battle0-elf", 29, 172, 15, `
#######       #######
#.G...#       #..E..#   E(158)
#...EG#       #...E.#   E(14)
#.#.#G#  -->  #.#.#.#
#..G#E#       #...#.#
#.....#       #.....#
#######       #######`),

		sbattle("battle1-elf", 33, 948, 4, `
#######       #######
#E..EG#       #.E.E.#   E(200), E(23)
#.#G.E#       #.#E..#   E(200)
#E.##E#  -->  #E.##E#   E(125), E(200)
#G..#.#       #.E.#.#   E(200)
#..E#.#       #...#.#
#######       #######`),

		sbattle("battle2-elf", 37, 94, 15, `
#######       #######
#E.G#.#       #.E.#.#   E(8)
#.#G..#       #.#E..#   E(86)
#G.#.G#  -->  #..#..#
#G..#.#       #...#.#
#...E.#       #.....#
#######       #######`),

		sbattle("battle3-elf", 39, 166, 12, `
#######       #######
#.E...#       #...E.#   E(14)
#.#..G#       #.#..E#   E(152)
#.###.#  -->  #.###.#
#E#G#G#       #.#.#.#
#...#G#       #...#.#
#######       #######`),

		sbattle("battle4-elf", 30, 38, 34, `
#########       #########   
#G......#       #.......#   
#.E.#...#       #.E.#...#   E(38)
#..##..G#       #..##...#   
#...##..#  -->  #...##..#   
#...#...#       #...#...#   
#.G...G.#       #.......#   
#.....G.#       #.......#   
#########       #########`),
	}

	buf := &bytes.Buffer{}
	for _, tm := range tests {
		gf, err := ParseGoblinFight(tm.start)
		if err != nil {
			t.Fatal(err)
		}

		setAttackStrength := func() {
			if tm.elfAttackStrength != 0 {
				gf.SetElfAttackStrength(tm.elfAttackStrength)
			}
		}

		setAttackStrength()
		step := 0
		for _, tt := range tm.t {
			for step < tt.after {
				gf.Simulate()
				step++
				buf.Reset()
				gf.Dump(buf, tm.withStat)
				t.Logf("%s after %d:\n%s\n", tm.label, step, buf.String())
			}
			got := clean15Layout(buf.String())
			want := clean15Layout(tt.layout)
			if got != want {
				t.Fatalf("%s after %d steps got:\n%s\n\nwant:\n%s\n", tm.label, tt.after, got, want)
			}
		}

		if tm.outcome == 0 {
			continue
		}

		gf, err = ParseGoblinFight(tm.start)
		if err != nil {
			t.Fatal(err)
		}
		setAttackStrength()
		got, want := findOutcome(t, gf, tm.label, true), tm.outcome
		if got != want {
			t.Fatalf("%s outcome is %d; want %d", tm.label, got, want)
		}

		gf, err = ParseGoblinFight(tm.start)
		if err != nil {
			t.Fatal(err)
		}
		setAttackStrength()
		got2 := gf.FindOutcome()
		if got != got2 {
			t.Fatalf("test findOutcome() != gf.FindOutcome(): %d != %d", got, got2)
		}
	}
}

func TestAoC15Puzzle(t *testing.T) {
	// puzzle
	gf, err := ParseGoblinFight(input15)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("1st puzzle outcome: %d", gf.FindOutcome())

	for strength := 4; strength <= 200; strength++ {
		gf, err := ParseGoblinFight(input15)
		if err != nil {
			t.Fatal(err)
		}

		ec := gf.ElfCount()

		gf.SetElfAttackStrength(strength)
		oc := gf.FindOutcome()
		if ec == gf.ElfCount() {
			t.Logf("2nd puzzle outcome: %d (strength: %d)", oc, strength)
			return
		}
	}

	t.Fatal("failed")
}

func findOutcome(t *testing.T, gf *GoblinFight, label string, withStat bool) int {
	step := 0
	buf := &bytes.Buffer{}
	for {
		advanced, fullturn := gf.Simulate()
		if !advanced {
			return -1
		}

		buf.Reset()
		gf.Dump(buf, withStat)
		t.Logf("%s after %d:\n%s\n", label, step, buf.String())

		hp, ok := gf.OutcomeHP()
		if ok {
			rounds := step
			if fullturn {
				rounds++
			}
			return rounds * hp
		}

		step++
	}
}

func clean15Layout(src string) string {
	v := strings.Split(src, "\n")
	i := 0
	for _, s := range v {
		s = strings.TrimSpace(s)
		if s != "" {
			v[i] = s
			i++
		}
	}
	return strings.Join(v[:i], "\n")
}
