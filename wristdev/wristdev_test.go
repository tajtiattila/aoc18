package wristdev

import (
	"fmt"
	"sort"
	"testing"
)

func TestWristDevSimple(t *testing.T) {
	a := Arch(4)
	s0 := a.State(3, 2, 1, 1)
	s1 := a.State(3, 2, 2, 1)

	var match []string
	for _, op := range Ops() {
		s := s0.Clone()
		s.Run(op, 2, 1, 2)
		if s.RegistersEqual(s1) {
			match = append(match, op.Name())
		}
	}

	sort.Strings(match)

	want := []string{"addi", "mulr", "seti"}

	if len(match) != len(want) {
		t.Fatalf("got %v; want %v", match, want)
	}

	for i := range want {
		if match[i] != want[i] {
			t.Fatalf("got %v; want %v", match, want)
		}
	}
}

func TestOpsSimple(t *testing.T) {
	type test struct {
		before *State

		op      Operator
		a, b, c int // arguments

		after *State
	}

	xop := func(before *State, opn string, a, b, c int, after *State) test {
		op := Op(opn)
		if op == nil {
			t.Fatalf("invalid op %s", opn)
		}
		return test{
			before: before,

			op: op,
			a:  a,
			b:  b,
			c:  c,

			after: after,
		}
	}

	a4 := Arch(4)

	tests := []test{
		xop(a4.State(3, 1, 2, 2), "eqrr", 1, 2, 3, a4.State(3, 1, 2, 0)),
		xop(a4.State(3, 1, 2, 2), "eqrr", 1, 2, 1, a4.State(3, 0, 2, 2)),
	}

	for _, tt := range tests {
		n := fmt.Sprintf("%s %d %d %d", tt.op.Name(), tt.a, tt.b, tt.c)
		t.Run(n, func(t *testing.T) {
			state := tt.before.Clone()
			state.Run(tt.op, tt.a, tt.b, tt.c)
			if !state.RegistersEqual(tt.after) {
				t.Fatalf("%s from %s got %s; want %s", n, tt.before.String(), state.String(), tt.after.String())
			}
		})
	}
}
