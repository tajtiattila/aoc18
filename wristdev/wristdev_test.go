package wristdev

import (
	"fmt"
	"sort"
	"testing"
)

func TestWristDev(t *testing.T) {
	s0 := St(3, 2, 1, 1)
	s1 := St(3, 2, 2, 1)

	var match []string
	for _, op := range Ops() {
		s := s0
		s.Run(op, 2, 1, 2)
		if s == s1 {
			match = append(match, op.Name())
		}
	}

	sort.Strings(match)

	want := []string{"addi", "mulr", "seti"}

	if len(match) != len(want) {
		t.Fatalf("got %v; want %v", (match), (want))
	}

	for i := range want {
		if match[i] != want[i] {
			t.Fatalf("got %v; want %v", (match), (want))
		}
	}
}

func TestOps(t *testing.T) {
	type test struct {
		before State

		op      Operator
		a, b, c int // arguments

		after State
	}

	xop := func(before State, opn string, a, b, c int, after State) test {
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

	tests := []test{
		xop(St(3, 1, 2, 2), "eqrr", 1, 2, 3, St(3, 1, 2, 0)),
		xop(St(3, 1, 2, 2), "eqrr", 1, 2, 1, St(3, 0, 2, 2)),
	}

	for _, tt := range tests {
		n := fmt.Sprintf("%s %d %d %d", tt.op.Name(), tt.a, tt.b, tt.c)
		t.Run(n, func(t *testing.T) {
			state := tt.before
			state.Run(tt.op, tt.a, tt.b, tt.c)
			if state != tt.after {
				t.Fatalf("%s from %s got %s; want %s", n, tt.before.String(), state.String(), tt.after.String())
			}
		})
	}
}
