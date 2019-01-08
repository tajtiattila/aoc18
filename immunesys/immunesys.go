package immunesys

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"sort"
	"strings"

	"github.com/pkg/errors"
)

type Group struct {
	Name      string
	UnitCount int // unit count
	HP        int // hp per unit

	Immune []string // attack types units immune to (no damage)
	Weak   []string // attack types units weak to (2x damage)

	Attack struct {
		Type   string // attack type
		Damage int    // damage per unit
	}

	Initiative int
}

func (g *Group) EffectivePower() int {
	return g.Attack.Damage * g.UnitCount
}

func (g *Group) AttackDamage(target *Group) int {
	if target == nil {
		return 0
	}

	for _, s := range target.Immune {
		if s == g.Attack.Type {
			return 0
		}
	}

	ep := g.EffectivePower()

	for _, s := range target.Weak {
		if s == g.Attack.Type {
			return 2 * ep
		}
	}

	return ep
}

func (g *Group) PerformAttack(target *Group) int {
	if target == nil {
		return 0
	}

	d := g.AttackDamage(target)

	kills := d / target.HP
	if kills > target.UnitCount {
		kills = target.UnitCount
	}

	target.UnitCount -= kills
	return kills
}

type Battle struct {
	Group map[string][]*Group
}

func ParseBattle(r io.Reader) (*Battle, error) {
	scanner := bufio.NewScanner(r)

	m := make(map[string][]*Group)
	var army string
	for lineno := 1; scanner.Scan(); lineno++ {
		line := scanner.Text()
		if line == "" {
			continue
		}

		if strings.HasSuffix(line, ":") {
			army = strings.TrimSuffix(line, ":")
			continue
		}

		if army == "" {
			return nil, errors.Errorf("no army on line %d", lineno)
		}

		const s1 = "hit points"
		const s2 = "with an attack"

		i1 := strings.Index(line, s1)
		i2 := strings.Index(line, s2)
		if i1 < 0 || i2 < 0 || i2 < i1 {
			return nil, errors.Errorf("invalid line %d", lineno)
		}

		i1 += len(s1)

		health, attrs, attack := line[:i1], line[i1:i2], line[i2:]

		g := &Group{}

		_, err := fmt.Sscanf(health, "%d units each with %d hit points",
			&g.UnitCount, &g.HP)
		if err != nil {
			return nil, errors.Wrapf(err, "line %d health", lineno)
		}

		if err := g.scanAttrs(attrs); err != nil {
			return nil, errors.Wrapf(err, "line %d attrs", lineno)
		}

		_, err = fmt.Sscanf(attack, "with an attack that does %d %s damage at initiative %d",
			&g.Attack.Damage, &g.Attack.Type, &g.Initiative)
		if err != nil {
			return nil, errors.Wrapf(err, "line %d attack", lineno)
		}

		i := len(m[army])
		g.Name = fmt.Sprintf("group %d", i+1)

		m[army] = append(m[army], g)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return &Battle{Group: m}, nil
}

func (g *Group) scanAttrs(a string) error {
	a = strings.TrimSpace(a)
	if a == "" {
		return nil
	}

	if a[0] == '(' && a[len(a)-1] == ')' {
		a = a[1 : len(a)-1]
	}

	for _, spec := range strings.Split(a, ";") {
		f := strings.Fields(spec)

		if len(f) < 2 || f[1] != "to" {
			return errors.Errorf("unrecognised attr specification: %s", spec)
		}

		for i := 2; i < len(f); i++ {
			f[i] = strings.TrimSuffix(f[i], ",")
		}

		switch f[0] {
		case "immune":
			g.Immune = append(g.Immune, f[2:]...)
		case "weak":
			g.Weak = append(g.Weak, f[2:]...)
		default:
			return errors.Errorf("unknown attr type: %q", f[0])
		}
	}

	return nil
}

func (b *Battle) ShowHeader(log io.Writer) {
	b.header(log, true)
}

type groupSpec struct {
	army  string
	group *Group

	target *Group
}

func (b *Battle) header(log io.Writer, showattack bool) []groupSpec {
	names := make([]string, 0, len(b.Group))
	for k := range b.Group {
		names = append(names, k)
	}
	sort.Strings(names)

	var group []groupSpec
	active := 0
	for _, k := range names {
		fmt.Fprintf(log, "%s:\n", k)
		for _, g := range b.Group[k] {
			var attack string
			if showattack {
				attack = fmt.Sprintf(" (damage %d)", g.Attack.Damage)
			}
			fmt.Fprintf(log, "%s contains %d units%s\n", strings.Title(g.Name), g.UnitCount, attack)
			group = append(group, groupSpec{
				army:  k,
				group: g,
			})
		}
		if len(b.Group[k]) == 0 {
			fmt.Fprintln(log, "No groups remain.")
		} else {
			active++
		}
	}

	if active < 2 {
		return nil
	}
	return group
}

func (b *Battle) Step(log io.Writer) bool {
	group := b.header(log, false)

	if len(group) == 0 {
		return false
	}

	fmt.Fprintln(log)

	starttc := b.TotalUnitCount()

	// target selection
	sort.Slice(group, func(i, j int) bool {
		gi := group[i].group
		gj := group[j].group

		ei := gi.EffectivePower()
		ej := gj.EffectivePower()
		if ei != ej {
			// higher effective power first
			return ei > ej
		}
		// otherwise initiative
		return gi.Initiative > gj.Initiative
	})

	var armyorder []string
	seenarmy := make(map[string]struct{})
	for _, g := range group {
		if _, ok := seenarmy[g.army]; !ok {
			armyorder = append(armyorder, g.army)
			seenarmy[g.army] = struct{}{}
		}
	}

	var target []groupSpec
	target = append(target, group...)

	for _, army := range armyorder {
		for i := range group {
			g := &group[i]
			if g.army != army {
				continue
			}
			bestDamage, bestIndex := -1, 0
			for i, t := range target {
				if t.group == nil || t.army == g.army {
					continue // already taken or own army
				}
				d := g.group.AttackDamage(t.group)
				fmt.Fprintf(log, "%s %s would deal defending %s %d damage\n",
					g.army, g.group.Name, t.group.Name, d)
				if d > bestDamage {
					bestDamage, bestIndex = d, i
				}
			}
			if bestDamage > 0 {
				g.target = target[bestIndex].group
				target[bestIndex].group = nil
			}
		}
	}

	fmt.Fprintln(log)

	// attacking
	sort.Slice(group, func(i, j int) bool {
		gi := group[i].group
		gj := group[j].group

		return gi.Initiative > gj.Initiative
	})

	for _, g := range group {
		if g.target == nil || g.target.UnitCount == 0 {
			continue
		}
		kills := g.group.PerformAttack(g.target)
		pl := "s"
		if kills == 1 {
			pl = ""
		}
		fmt.Fprintf(log, "%s %s attacks defending %s, killing %d unit%s\n",
			g.army, g.group.Name, g.target.Name, kills, pl)
	}

	fmt.Fprintln(log)

	// remove killed groups
	for k := range b.Group {
		v := b.Group[k]
		j := 0
		for _, g := range v {
			if g.UnitCount > 0 {
				v[j] = g
				j++
			}
		}
		b.Group[k] = v[:j]
	}

	return starttc != b.TotalUnitCount()
}

func (b *Battle) TotalUnitCount() int {
	n := 0
	for _, v := range b.Group {
		for _, g := range v {
			n += g.UnitCount
		}
	}
	return n
}

func (b *Battle) Clone() *Battle {
	m := make(map[string][]*Group)
	for k, v := range b.Group {
		w := make([]*Group, len(v))
		for i := range v {
			g := *v[i]
			w[i] = &g
		}
		m[k] = w
	}
	return &Battle{
		Group: m,
	}
}

func (b *Battle) Boost(army string, attackBoost int) {
	v, ok := b.Group[army]
	if !ok {
		panic("invalid boost")
	}

	for _, g := range v {
		g.Attack.Damage += attackBoost
	}
}

func (b *Battle) Run() (winner string, finished bool) {
	for b.Step(ioutil.Discard) {
	}
	n := 0
	for k, v := range b.Group {
		if len(v) > 0 {
			winner = k
			n++
		}
	}

	if n == 1 {
		return winner, true
	}
	return "", false
}
