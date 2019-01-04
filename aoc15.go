package main

import (
	"bytes"
	"fmt"
	"io"
	"strings"

	"github.com/pkg/errors"
)

type GoblinFight struct {
	dx, dy int

	m []gftile

	fm []int // map of path distances

	// last flood info
	lastf struct {
		ix, ax int
		iy, ay int
	}

	elf, goblin gfteam
}

type gfteam struct {
	teamHP int
	count  int

	attackStrength int
}

// fill value for GoblinFight.fm
const gfmEmpty = 1e7
const gfmBlock = 1e6

func ParseGoblinFight(layout string) (*GoblinFight, error) {
	src := strings.Split(strings.TrimSpace(layout), "\n")
	dy := len(src)
	if dy == 0 {
		return nil, errors.New("empty source")
	}

	dx := len(src[0])
	if dx == 0 {
		return nil, errors.New("empty source")
	}

	gf := &GoblinFight{
		dx: dx,
		dy: dy,
		m:  make([]gftile, dx*dy),
		fm: make([]int, dx*dy),

		elf: gfteam{
			attackStrength: gfDefaultAttackStrength,
		},
		goblin: gfteam{
			attackStrength: gfDefaultAttackStrength,
		},
	}

	for i := range gf.fm {
		gf.fm[i] = gfmEmpty
	}

	runetile := func(r rune) gftile {
		var k, hp gftile
		switch r {
		case '.':
			k = gfSpace
		case '#':
			k = gfWall
		case 'E':
			k, hp = gfElf, gfStartHP
			gf.elf.teamHP += gfStartHP
			gf.elf.count++
		case 'G':
			k, hp = gfGoblin, gfStartHP
			gf.goblin.teamHP += gfStartHP
			gf.elf.count++
		}

		return (k << gfKindShift) | hp
	}

	for y, line := range src {
		if len(line) != dx {
			return nil, errors.Errorf("line %d invalid length", y)
		}

		ofs := gf.ofs(0, y)
		for _, r := range line {
			t := runetile(r)
			if t == 0 && r != '.' {
				return nil, errors.Errorf("line %d has invalid rune %c", y, r)
			}
			gf.m[ofs] = t
			ofs++
		}
	}

	return gf, nil
}

func (gf *GoblinFight) ElfCount() int {
	return gf.elf.count
}

func (gf *GoblinFight) SetElfAttackStrength(v int) {
	gf.elf.attackStrength = v
}

func (gf *GoblinFight) ElvesWon() bool {
	return gf.goblin.teamHP == 0 && gf.elf.teamHP > 0
}

type gftile uint16

func gftileKindHP(kind, hp int) gftile {
	k := gftile(kind)
	h := gftile(hp)
	return (k << gfKindShift) | h
}

func (t gftile) Kind() int { return int(uint16(t) >> 12) }
func (t gftile) HP() int   { return int(uint16(t) & gfHPMask) }

const (
	gfSpace = iota
	gfWall
	gfElf
	gfGoblin

	gfKindShift = 12
	gfHPMask    = (1 << gfKindShift) - 1

	gfStartHP               = 200
	gfDefaultAttackStrength = 3
)

func (gf *GoblinFight) Simulate() (advanced, fullturn bool) {
	skip := make(map[Point]struct{})
	advanced, fullturn = false, true
	for y := 0; y < gf.dy; y++ {
		for x := 0; x < gf.dx; x++ {
			p := Pt(x, y)
			if _, ok := skip[p]; ok {
				continue
			}

			k := gf.m[gf.pofs(p)].Kind()

			if k == gfElf || k == gfGoblin {
				if gf.elf.teamHP*gf.goblin.teamHP == 0 {
					fullturn = false
				}

				goal, _ := gf.pickGoal(p)
				if goal != p {
					p = gf.step(p, goal)
					skip[p] = struct{}{}
					advanced = true
				}
				if gf.attack(p, skip) {
					advanced = true
				}
			}
		}
	}
	return advanced, fullturn
}

func (gf *GoblinFight) FindOutcome() int {
	step := 0
	for {
		advanced, fullturn := gf.Simulate()
		if !advanced {
			return -1
		}

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

func (gf *GoblinFight) OutcomeHP() (outcome int, ok bool) {
	var e, g int
	for _, t := range gf.m {
		switch t.Kind() {
		case gfElf:
			if g != 0 {
				return 0, false
			}
			e += t.HP()
		case gfGoblin:
			if e != 0 {
				return 0, false
			}
			g += t.HP()
		}
	}

	return e + g, true
}

func (gf *GoblinFight) ofs(x, y int) int { return x + y*gf.dx }
func (gf *GoblinFight) pofs(p Point) int { return p.X + p.Y*gf.dx }

func (gf *GoblinFight) tile(x, y int) gftile { return gf.m[gf.ofs(x, y)] }

// pickGoal picks a goal (position to move to) for the unit at p.
func (gf *GoblinFight) pickGoal(p Point) (goal Point, goalDist int) {

	var enemyKind int
	switch gf.tile(p.X, p.Y).Kind() {
	case gfElf:
		enemyKind = gfGoblin
	case gfGoblin:
		enemyKind = gfElf
	default:
		panic("only units can pick goal")
	}

	var enemy []Point

	enemyDist := gf.flood(p.X, p.Y, func(next Point) gflResult {
		switch gf.tile(next.X, next.Y).Kind() {
		case enemyKind:
			enemy = append(enemy, next)
			return gflStop
		case gfSpace:
			return gflContinue
		default:
			return gflBlock
		}
	})

	if len(enemy) == 0 {
		return p, -1
	}

	// distance of target position
	goalDist = enemyDist - 1
	if goalDist == 0 {
		// already next to an enemy
		return p, goalDist
	}

	first := true
	var adjbuf [4]Point
	for _, enemyp := range enemy {
		// check points next to an enemy
		for _, q := range gf.adj(enemyp, adjbuf[:0]) {
			if gf.fm[gf.pofs(q)] == goalDist {
				if first || q.Y < goal.Y ||
					(q.Y == goal.Y && q.X < goal.X) {
					goal = q
					first = false
				}
			}
		}
	}
	return goal, goalDist
}

// step walks the unit at p one step along goal,
// and returns the new position.
func (gf *GoblinFight) step(p, goal Point) Point {
	if p == goal {
		return p
	}

	unit := gf.m[gf.pofs(p)]
	if unit.Kind() != gfElf && unit.Kind() != gfGoblin {
		panic("only units can step")
	}

	v := gf.flood(goal.X, goal.Y, func(w Point) gflResult {
		switch {
		case w == p:
			return gflStop
		case gf.m[gf.pofs(w)].Kind() != gfSpace:
			return gflBlock
		default:
			return gflContinue
		}
	})
	v--

	var next Point
	var adjbuf [4]Point
	for _, next = range gf.adj(p, adjbuf[:0]) {
		if gf.fm[gf.pofs(next)] == v {
			break
		}
	}

	gf.m[gf.pofs(p)] = gfSpace
	gf.m[gf.pofs(next)] = unit
	return next
}

func (gf *GoblinFight) attack(p Point, skip map[Point]struct{}) bool {
	var enemyKind, attackStrength int
	var enemyTeam *gfteam
	switch gf.tile(p.X, p.Y).Kind() {
	case gfElf:
		attackStrength = gf.elf.attackStrength
		enemyKind = gfGoblin
		enemyTeam = &gf.goblin
	case gfGoblin:
		attackStrength = gf.goblin.attackStrength
		enemyKind = gfElf
		enemyTeam = &gf.elf
	default:
		panic("only units can attack")
	}

	targetHP := 0
	var target Point
	var adjbuf [4]Point
	for _, q := range gf.adj(p, adjbuf[:0]) {
		t := gf.tile(q.X, q.Y)
		if t.Kind() == enemyKind && (targetHP == 0 || t.HP() < targetHP) {
			target, targetHP = q, t.HP()
		}
	}

	if targetHP == 0 {
		// nothing to attack
		return false
	}

	var t gftile
	if targetHP <= attackStrength {
		// killed
		t = gfSpace
		enemyTeam.teamHP -= targetHP
		enemyTeam.count--
		delete(skip, target)
	} else {
		targetHP -= attackStrength
		enemyTeam.teamHP -= attackStrength
		t = gftileKindHP(enemyKind, targetHP)
	}
	gf.m[gf.ofs(target.X, target.Y)] = t
	return true
}

func (gf *GoblinFight) flood(sx, sy int, floodf func(p Point) gflResult) int {
	// clear last flood info
	cofs := gf.ofs(gf.lastf.ix, gf.lastf.iy)
	for y := gf.lastf.iy; y <= gf.lastf.ay; y++ {
		ofs := cofs
		cofs += gf.dx
		for x := gf.lastf.ix; x <= gf.lastf.ax; x++ {
			gf.fm[ofs] = gfmEmpty
			ofs++
		}
	}

	gf.lastf.ix = sx
	gf.lastf.iy = sy
	gf.lastf.ax = sx
	gf.lastf.ay = sy

	dist := 0
	active := []Point{Pt(sx, sy)}
	ofs := gf.ofs(sx, sy)
	gf.fm[ofs] = dist
	var nactive []Point
	var adjbuf [4]Point
	done := false
	for !done && len(active) != 0 {
		dist++
		for _, q := range active {
			for _, p := range gf.adj(q, adjbuf[:0]) {
				ofs := gf.pofs(p)
				if gf.fm[ofs] != gfmEmpty {
					continue // processed
				}

				// extend flood bounding box
				if p.X < gf.lastf.ix {
					gf.lastf.ix = p.X
				}
				if p.Y < gf.lastf.iy {
					gf.lastf.iy = p.Y
				}
				if p.X > gf.lastf.ax {
					gf.lastf.ax = p.X
				}
				if p.Y > gf.lastf.ay {
					gf.lastf.ay = p.Y
				}

				switch floodf(p) {

				case gflBlock:
					gf.fm[ofs] = gfmBlock

				case gflContinue:
					gf.fm[ofs] = dist
					nactive = append(nactive, p)

				case gflStop:
					gf.fm[ofs] = dist
					done = true
				}
			}
		}

		active, nactive = nactive, active[:0]
	}

	return dist
}

func (gf *GoblinFight) adj(p Point, res []Point) []Point {
	if p.Y > 0 {
		res = append(res, Pt(p.X, p.Y-1))
	}
	if p.X > 0 {
		res = append(res, Pt(p.X-1, p.Y))
	}
	if p.X+1 < gf.dx {
		res = append(res, Pt(p.X+1, p.Y))
	}
	if p.Y+1 < gf.dy {
		res = append(res, Pt(p.X, p.Y+1))
	}
	return res
}

type gflResult int

const (
	gflBlock gflResult = iota // path blocked
	gflContinue
	gflStop // finish after processing tiles with the same distance
)

func (gf *GoblinFight) Dump(w io.Writer, withStat bool) {

	var buf bytes.Buffer
	var addStat func(r rune, hp int)
	if withStat {
		addStat = func(r rune, hp int) {
			if buf.Len() == 0 {
				buf.WriteString("   ")
			} else {
				buf.WriteString(", ")
			}
			fmt.Fprintf(&buf, "%c(%d)", r, hp)
		}
	} else {
		addStat = func(rune, int) {}
	}

	for i, t := range gf.m {
		var r rune
		switch t.Kind() {
		case gfSpace:
			r = '.'
		case gfWall:
			r = '#'
		case gfElf:
			r = 'E'
			addStat(r, t.HP())
		case gfGoblin:
			r = 'G'
			addStat(r, t.HP())
		default:
			r = '?'
		}
		fmt.Fprintf(w, "%c", r)

		if i%gf.dx == gf.dx-1 {
			fmt.Fprint(w, buf.String(), "\n")
			buf.Reset()
		}
	}
}

func (gf *GoblinFight) DumpFloodMap(w io.Writer) {
	for i, v := range gf.fm {
		switch v {
		case gfmEmpty:
			fmt.Fprint(w, "  .")
		case gfmBlock:
			fmt.Fprint(w, "  #")
		default:
			fmt.Fprintf(w, "%3d", v)
		}
		if i%gf.dx == gf.dx-1 {
			fmt.Fprint(w, "\n")
		}
	}
}
