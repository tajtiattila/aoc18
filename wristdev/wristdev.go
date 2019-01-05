package wristdev

import (
	"bytes"
	"fmt"

	"github.com/pkg/errors"
)

type Architecture struct {
	NReg int // number of registers

	// instruction pointer register
	IP struct {
		IsRegister bool
		Index      int
	}
}

func Arch(numRegisters int) *Architecture {
	return &Architecture{
		NReg: numRegisters,
	}
}

func ArchWithIP(numRegisters, ipIndex int) *Architecture {
	if ipIndex >= numRegisters {
		panic("invalid IP architecture")
	}
	a := &Architecture{
		NReg: numRegisters,
	}
	a.IP.IsRegister = true
	a.IP.Index = ipIndex
	return a
}

// State creates a new empty state for a
// from the values specified.
func (a *Architecture) State(values ...int) *State {
	if a.NReg <= 0 {
		panic("invalid register architecture")
	}
	if a.IP.IsRegister && (a.IP.Index < 0 && a.IP.Index >= a.NReg) {
		panic("invalid IP architecture")
	}

	s := &State{
		Arch: a,
		R:    make([]int, a.NReg),
	}
	if a.IP.IsRegister {
		s.IP = &s.R[a.IP.Index]
	} else {
		s.IP = new(int)
	}

	copy(s.R, values)

	return s
}

type State struct {
	Arch *Architecture

	R  []int // registers
	IP *int
}

func (old *State) Clone() *State {
	a := old.Arch

	s := &State{
		Arch: a,
		R:    make([]int, a.NReg),
	}

	if a.IP.IsRegister {
		s.IP = &s.R[a.IP.Index]
	} else {
		s.IP = new(int)
		*s.IP = *old.IP
	}

	copy(s.R, old.R)

	return s
}

func (s *State) RegistersEqual(r *State) bool {
	if len(s.R) != len(r.R) {
		return false
	}

	for i := range s.R {
		if s.R[i] != r.R[i] {
			return false
		}
	}
	return true
}

func (s *State) String() string {
	var pfx string
	if s.Arch.IP.IsRegister {
		pfx = fmt.Sprintf("ip=%d ", *s.IP)
	}
	return pfx + fmt.Sprint(s.R)
}

func (s *State) Run(op Operator, a, b, c int) {
	op.Run(s, a, b, c)
	*(s.IP)++
}

// Step steps once in prog, and reports if the process is still running.
func (s *State) Step(prog []Instruction) bool {
	if *s.IP < 0 || *s.IP >= len(prog) {
		return false
	}

	inst := prog[*s.IP]
	s.Run(inst.Op, inst.A, inst.B, inst.C)

	return 0 <= *s.IP && *s.IP < len(prog)
}

func (s *State) RunProgram(prog []Instruction, maxstep int) bool {
	if *s.IP < 0 || *s.IP >= len(prog) {
		return false
	}

	for i := 0; i < maxstep && 0 <= *s.IP && *s.IP < len(prog); i++ {
		inst := prog[*s.IP]
		s.Run(inst.Op, inst.A, inst.B, inst.C)
	}

	return 0 <= *s.IP && *s.IP < len(prog)
}

func (s State) reg(n int) int {
	if 0 <= n && n < len(s.R) {
		return s.R[n]
	}
	return 0
}

func (s *State) preg(n int) *int {
	if 0 <= n && n < len(s.R) {
		return &s.R[n]
	}
	return nil
}

type Operator interface {
	Name() string // full mnemonic

	// Run runs the instruction on s
	Run(s *State, a, b, c int)

	Prefix() string // mnemonic prefix, eg. "add" for "addi"

	// Args reports if argument a and/or b immediate
	Args() (a, b ArgKind)
}

type ArgKind int

const (
	ArgUnused ArgKind = iota // unused argument (seti/setr)
	ArgImmediate
	ArgReg
)

type opSimple struct {
	prefix string
	name   string

	at, bt ArgKind

	f func(s *State, a, b int) int
}

func (o opSimple) Prefix() string       { return o.prefix }
func (o opSimple) Name() string         { return o.name }
func (o opSimple) Args() (a, b ArgKind) { return o.at, o.bt }

func (o opSimple) Run(s *State, a, b, c int) {
	r := o.f(s, a, b)
	if p := s.preg(c); p != nil {
		*p = r
	}
}

type opImmediate struct {
	prefix string
	name   string

	f func(a, b int) int
}

func (o opImmediate) Prefix() string       { return o.prefix }
func (o opImmediate) Name() string         { return o.name }
func (o opImmediate) Args() (a, b ArgKind) { return ArgReg, ArgImmediate }

func (o opImmediate) Run(s *State, a, b, c int) {
	r := o.f(s.reg(a), b)
	if p := s.preg(c); p != nil {
		*p = r
	}
}

type opReg struct {
	prefix string
	name   string

	f func(a, b int) int
}

func (o opReg) Prefix() string       { return o.prefix }
func (o opReg) Name() string         { return o.name }
func (o opReg) Args() (a, b ArgKind) { return ArgReg, ArgReg }

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
				prefix: prefix,
				name:   prefix + "i",
				f:      f,
			},
			opReg{
				prefix: prefix,
				name:   prefix + "r",
				f:      f,
			})
	}

	type argkindpair struct {
		a, b ArgKind
	}
	akp := func(a, b ArgKind) argkindpair {
		return argkindpair{a: a, b: b}
	}
	ms := map[string]argkindpair{
		"r":  akp(ArgReg, ArgUnused),
		"i":  akp(ArgImmediate, ArgUnused),
		"ir": akp(ArgImmediate, ArgReg),
		"rr": akp(ArgReg, ArgReg),
		"ri": akp(ArgReg, ArgImmediate),
	}

	adds := func(prefix, suffix string, f func(s *State, a, b int) int) {
		akp, ok := ms[suffix]
		if !ok {
			panic(fmt.Sprintf("invalid suffix: %q", suffix))
		}
		ops = append(ops,
			opSimple{
				prefix: prefix,
				name:   prefix + suffix,
				at:     akp.a,
				bt:     akp.b,
				f:      f,
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

	adds("set", "i", func(s *State, a, b int) int {
		return a
	})
	adds("set", "r", func(s *State, a, b int) int {
		return s.reg(a)
	})

	adds("gt", "ir", func(s *State, a, b int) int {
		if a > s.reg(b) {
			return 1
		}
		return 0
	})
	adds("gt", "ri", func(s *State, a, b int) int {
		if s.reg(a) > b {
			return 1
		}
		return 0
	})
	adds("gt", "rr", func(s *State, a, b int) int {
		if s.reg(a) > s.reg(b) {
			return 1
		}
		return 0
	})

	adds("eq", "ir", func(s *State, a, b int) int {
		if a == s.reg(b) {
			return 1
		}
		return 0
	})
	adds("eq", "ri", func(s *State, a, b int) int {
		if s.reg(a) == b {
			return 1
		}
		return 0
	})
	adds("eq", "rr", func(s *State, a, b int) int {
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

type Instruction struct {
	Op Operator

	A, B, C int
}

func ParseInstruction(s string) (Instruction, error) {
	var name string
	var a, b, c int
	_, err := fmt.Sscanf(s, "%s %d %d %d", &name, &a, &b, &c)
	if err != nil {
		return Instruction{}, errors.Errorf("ParseInstruction: can't parse %q", s)
	}

	op := Op(name)
	if op == nil {
		return Instruction{}, errors.Errorf("ParseInstruction: op %q unknown", name)
	}

	return Instruction{
		Op: op,
		A:  a,
		B:  b,
		C:  c,
	}, nil
}

func ParseProgram(v []string) ([]Instruction, error) {
	var prog []Instruction
	for no, line := range v {
		inst, err := ParseInstruction(line)
		if err != nil {
			return prog, errors.Wrapf(err, "line %d", no+1)
		}
		prog = append(prog, inst)
	}
	return prog, nil
}

func (i Instruction) String() string {
	return fmt.Sprintf("%s %d %d %d", i.Op.Name(), i.A, i.B, i.C)
}

func (i Instruction) Fmt(regnames []string) string {
	var buf bytes.Buffer
	buf.WriteString(i.Op.Prefix())

	ak, bk := i.Op.Args()

	s := argName(regnames, ak, i.A)
	if s != "" {
		buf.WriteString(" ")
		buf.WriteString(s)
	}

	s = argName(regnames, bk, i.B)
	if s != "" {
		buf.WriteString(" ")
		buf.WriteString(s)
	}

	s = argName(regnames, ArgReg, i.C)
	buf.WriteString(" ")
	buf.WriteString(s)

	return buf.String()
}

func argName(regnames []string, k ArgKind, v int) string {
	switch k {
	case ArgUnused:
		return ""
	case ArgImmediate:
		return fmt.Sprintf("%d", v)
	case ArgReg:
		if v < len(regnames) {
			return regnames[v]
		}
		return fmt.Sprintf("r%d", v)
	default:
		panic("invalid ArgKind")
	}
}
