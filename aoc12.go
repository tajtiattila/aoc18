package main

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/pkg/errors"
)

type CavePots struct {
	zero int // offset of pot zero in slice

	p []byte // pot data; 0: no plant, 1: plant

	last []byte // last pot data

	rule []cavePotRule

	maxstepgrow int
}

func (cp *CavePots) clone() *CavePots {
	return &CavePots{
		zero:        cp.zero,
		p:           append([]byte(nil), cp.p...),
		rule:        cp.rule,
		maxstepgrow: cp.maxstepgrow,
	}
}

func NewCavePots(src string) (*CavePots, error) {
	lines := strings.Split(src, "\n")
	if len(lines) == 0 {
		return nil, errors.New("missing initial state")
	}

	var is string
	if _, err := fmt.Sscanf(lines[0], "initial state: %s", &is); err != nil {
		return nil, errors.Wrap(err, "cannot scan initial state")
	}

	if !potStateValid(is) {
		return nil, errors.New("invalid character in initial state")
	}

	cp := &CavePots{
		p: []byte(is),
	}

	for _, line := range lines[1:] {
		if line == "" {
			continue
		}

		var cs, rs string
		if _, err := fmt.Sscanf(line, "%s => %s", &cs, &rs); err != nil {
			return nil, errors.Wrapf(err, "cannot scan rule line %q initial state", line)
		}

		if !potStateValid(cs) || !potStateValid(rs) {
			return nil, errors.Errorf("invalid character in rule %q", line)
		}

		if len(cs)%2 != 1 || len(rs) != 1 {
			return nil, errors.Errorf("invalid length in rule %q", line)
		}
		cp.rule = append(cp.rule, cavePotRule{
			cond:   []byte(cs),
			result: rs[0],
		})
		if n := len(cs); n > cp.maxstepgrow {
			cp.maxstepgrow = n
		}
	}

	return cp, nil
}

func potStateValid(s string) bool {
	for _, r := range s {
		switch r {
		case '.', '#':
			// pass
		default:
			return false
		}
	}
	return true
}

func (cp *CavePots) firstPlantIdx() int {
	for i, v := range cp.p {
		if v == '#' {
			return i
		}
	}
	return len(cp.p)
}

func (cp *CavePots) addBorder(n int) {
	cp.p = append(potSlice(n), cp.p...)
	cp.p = append(cp.p, potSlice(n)...)
	cp.zero += n
}

func (cp *CavePots) ensureBorder() {
	n := cp.maxstepgrow
	if len(cp.p) < 2*n {
		cp.addBorder(n)
		return
	}

	grow := len(cp.p)
	if havePlant(cp.p[:n]) {
		cp.p = append(potSlice(grow), cp.p...)
		cp.zero += grow
	}
	if havePlant(cp.p[len(cp.p)-n:]) {
		cp.p = append(cp.p, potSlice(grow)...)
	}
}

func potSlice(n int) []byte {
	p := make([]byte, n)
	for i := range p {
		p[i] = '.'
	}
	return p
}

func havePlant(p []byte) bool {
	for _, c := range p {
		if c == '#' {
			return true
		}
	}
	return false
}

func (cp *CavePots) Step() {
	cp.ensureBorder()

	if len(cp.last) < len(cp.p) {
		cp.last = make([]byte, len(cp.p))
	}
	copy(cp.last, cp.p)

	n := cp.maxstepgrow / 2
	for i := n; i < len(cp.last)-n; i++ {
		cp.stepAt(i)
	}
}

type cavePotRule struct {
	cond []byte // odd length, middle item is current

	result byte
}

func (cp *CavePots) stepAt(i int) {
	for _, rule := range cp.rule {
		lc := len(rule.cond)
		ofs := i - lc/2

		if bytes.Equal(rule.cond, cp.last[ofs:ofs+lc]) {
			cp.p[i] = rule.result
			return
		}
	}
	cp.p[i] = '.'
}

func (cp *CavePots) NumPlants() int {
	n := 0
	for _, c := range cp.p {
		if c == '#' {
			n++
		}
	}
	return n
}

func (cp *CavePots) PlantSum() int {
	n := 0
	for i, c := range cp.p {
		if c == '#' {
			index := i - cp.zero
			n += index
		}
	}
	return n
}

func (cp *CavePots) Fmt(start, width int) string {
	i := cp.zero + start
	wsave := width

	buf := &bytes.Buffer{}
	if i < 0 {
		for i < 0 {
			buf.WriteByte('.')
			i++
			width--
		}
	}

	if i+width < len(cp.p) {
		buf.Write(cp.p[i : i+width])
	} else {
		buf.Write(cp.p[i:])
		for buf.Len() < wsave {
			buf.WriteByte('.')
		}
	}

	return buf.String()
}
