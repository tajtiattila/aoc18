package wristdev

import "fmt"

const NRegisters = 4

type State [NRegisters]int

func St(a, b, c, d int) State {
	return State{a, b, c, d}
}

func (s State) String() string {
	return "[" + fmt.Sprint(s[:]) + "]"
}

func (s *State) Run(op Operator, a, b, c int) {
	op.Run(s, a, b, c)
}

func (s State) reg(n int) int {
	if 0 <= n && n < NRegisters {
		return s[n]
	}
	return 0
}

func (s *State) preg(n int) *int {
	if 0 <= n && n < NRegisters {
		return &s[n]
	}
	return nil
}

type Operator interface {
	Name() string

	// Run runs the instruction on s
	Run(s *State, a, b, c int)
}

type opSimple struct {
	name string

	f func(s *State, a, b int) int
}

func (i opSimple) Name() string { return i.name }

func (i opSimple) Run(s *State, a, b, c int) {
	r := i.f(s, a, b)
	if p := s.preg(c); p != nil {
		*p = r
	}
}

type opImmediate struct {
	name string

	f func(a, b int) int
}

func (o opImmediate) Name() string { return o.name }

func (o opImmediate) Run(s *State, a, b, c int) {
	r := o.f(s.reg(a), b)
	if p := s.preg(c); p != nil {
		*p = r
	}
}

type opReg struct {
	name string

	f func(a, b int) int
}

func (o opReg) Name() string { return o.name }

func (o opReg) Run(s *State, a, b, c int) {
	r := o.f(s.reg(a), s.reg(b))
	if p := s.preg(c); p != nil {
		*p = r
	}
}

func Op(name string) Operator {
	return opmap[name]
}

func Ops() []Operator {
	return ops
}

var (
	ops   []Operator
	opmap map[string]Operator
)

func init() {

	add := func(prefix string, f func(a, b int) int) {
		ops = append(ops,
			opImmediate{
				name: prefix + "i",
				f:    f,
			},
			opReg{
				name: prefix + "r",
				f:    f,
			})
	}

	adds := func(name string, f func(s *State, a, b int) int) {
		ops = append(ops,
			opSimple{
				name: name,
				f:    f,
			})
	}

	add("add", func(a, b int) int {
		return a + b
	})
	add("mul", func(a, b int) int {
		return a * b
	})
	add("ban", func(a, b int) int {
		return a & b
	})
	add("bor", func(a, b int) int {
		return a | b
	})

	adds("seti", func(s *State, a, b int) int {
		return a
	})
	adds("setr", func(s *State, a, b int) int {
		return s.reg(a)
	})

	adds("gtir", func(s *State, a, b int) int {
		if a > s.reg(b) {
			return 1
		}
		return 0
	})
	adds("gtri", func(s *State, a, b int) int {
		if s.reg(a) > b {
			return 1
		}
		return 0
	})
	adds("gtrr", func(s *State, a, b int) int {
		if s.reg(a) > s.reg(b) {
			return 1
		}
		return 0
	})

	adds("eqir", func(s *State, a, b int) int {
		if a == s.reg(b) {
			return 1
		}
		return 0
	})
	adds("eqri", func(s *State, a, b int) int {
		if s.reg(a) == b {
			return 1
		}
		return 0
	})
	adds("eqrr", func(s *State, a, b int) int {
		if s.reg(a) == s.reg(b) {
			return 1
		}
		return 0
	})

	opmap = make(map[string]Operator)
	for _, o := range ops {
		opmap[o.Name()] = o
	}
}
