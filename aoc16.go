package main

import (
	"fmt"
	"log"

	"github.com/tajtiattila/aoc18/wristdev"
)

func wristdevhack() {
	samples, program := puzzleInput16()

	n := 0
	for _, s := range samples {
		count := s.countMatchingOp()
		if count >= 3 {
			n++
		}
	}
	fmt.Println("16/1:", n)

	codeleft := make(map[int]struct{})
	for _, s := range samples {
		codeleft[s.instr.opcode] = struct{}{}
	}

	// init map with names of all ops
	nameleft := make(map[string]struct{})
	for _, op := range wristdev.Ops() {
		nameleft[op.Name()] = struct{}{}
	}

	codeop := make(map[int]wristdev.Operator)
	for len(codeleft) > 0 {

		if len(nameleft) == 0 {
			log.Fatal("all names taken")
		}

		for opcode := range codeleft {

			possible := make(map[string]struct{})
			for n := range nameleft {
				possible[n] = struct{}{}
			}

		SampleLoop:
			for _, s := range samples {
				if s.instr.opcode != opcode {
					continue
				}

				canbe := make(map[string]struct{})
				for _, op := range s.possibleOps(nil) {
					canbe[op.Name()] = struct{}{}
				}

				// remove impossible ops
				for name := range possible {
					if _, ok := canbe[name]; !ok {
						delete(possible, name)
					}
				}

				switch len(possible) {
				case 0:
					log.Fatalf("no possible op left for code %d", opcode)
				case 1:
					break SampleLoop
				}
			}

			if len(possible) == 1 {
				var name string
				for name = range possible {
				}

				if verbose {
					fmt.Printf("opcode %d is %s\n", opcode, name)
				}

				codeop[opcode] = wristdev.Op(name)
				delete(codeleft, opcode)
				delete(nameleft, name)
			}
		}
	}

	var state wristdev.State

	for _, instr := range program {
		op := codeop[instr.opcode]
		a, b, c := instr.args()
		state.Run(op, a, b, c)
	}

	fmt.Println("16/2:", state[0])
}

type wristdevInstr struct {
	opcode, a, b, c int
}

func (i wristdevInstr) args() (a, b, c int) {
	return i.a, i.b, i.c
}

type wristdevSample struct {
	before wristdev.State

	instr wristdevInstr

	after wristdev.State
}

func (s *wristdevSample) possibleOps(dst []wristdev.Operator) []wristdev.Operator {
	for _, op := range wristdev.Ops() {
		state := s.before
		a, b, c := s.instr.args()
		state.Run(op, a, b, c)
		if state == s.after {
			dst = append(dst, op)
		}
	}
	return dst
}

func (s *wristdevSample) countMatchingOp() int {
	return len(s.possibleOps(nil))
}

func puzzleInput16() (samples []wristdevSample, program []wristdevInstr) {
	lines := PuzzleInputLines(16)

	i := 0
	for i+3 < len(lines) {
		if lines[i] == "" {
			i++
			continue
		}

		var a, b, c, d int
		_, err := fmt.Sscanf(lines[i], "Before: [%d, %d, %d, %d]", &a, &b, &c, &d)
		if err != nil {
			break
		}

		before := wristdev.St(a, b, c, d)

		instr, err := parsewristdevInstr(lines[i+1])
		if err != nil {
			log.Fatalf("error parsing instruciton near line %d\n", i+2)
		}

		_, err = fmt.Sscanf(lines[i+2], "After:  [%d, %d, %d, %d]", &a, &b, &c, &d)
		if err != nil {
			log.Fatalf("error parsing state near line %d\n", i+3)
		}

		after := wristdev.St(a, b, c, d)

		samples = append(samples, wristdevSample{
			before: before,
			instr:  instr,
			after:  after,
		})

		i += 3
	}

	for ; i < len(lines); i++ {
		if lines[i] != "" {
			instr, err := parsewristdevInstr(lines[i])
			if err != nil {
				log.Fatalf("error parsing instruciton near line %d\n", i+1)
			}
			program = append(program, instr)
		}
	}

	return samples, program
}

func parsewristdevInstr(line string) (wristdevInstr, error) {
	var op, a, b, c int
	_, err := fmt.Sscanf(line, "%d%d%d%d", &op, &a, &b, &c)
	if err != nil {
		return wristdevInstr{}, err
	}
	return wristdevInstr{
		opcode: op,

		a: a,
		b: b,
		c: c,
	}, nil
}
