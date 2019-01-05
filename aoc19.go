package main

import (
	"fmt"
	"log"

	"github.com/tajtiattila/aoc18/wristdev"
)

func wristdev19() {
	ipreg, prog := puzzleInput19()

	const nreg = 6

	regnames := make([]string, nreg)
	for i := range regnames {
		regnames[i] = fmt.Sprintf("r%d", i)
	}
	regnames[ipreg] = "ip"

	if verbose {
		fmt.Println("disassembly:")
		for i, l := range prog {
			fmt.Printf("%3d  %s\n", i, l.Fmt(regnames))
		}
		fmt.Println()
	}

	arch := wristdev.ArchWithIP(nreg, ipreg)
	state := arch.State()

	runwristprog(state, prog)

	fmt.Println("19/1:", state.R[0])
	r0, r3 := aoc19sim(0)
	if verbose {
		fmt.Println(" Simulated:", r0, r3)
		fmt.Println(" Divsum:", aoc19divSum(r3))

		showr3 := func(r0 int) {
			state := arch.State(r0)
			for *state.IP != 1 {
				state.Step(prog)
			}
			fmt.Printf("r0=%d -> r3=%d\n", r0, state.R[3])
		}

		showr3(0)
		showr3(1)
	}

	r3 = aoc19init(1)
	fmt.Println("19/2:", aoc19divSum(r3))
}

func runwristprog(state *wristdev.State, prog []wristdev.Instruction) bool {

	w := verbosew

	if !verbose {
		state.RunProgram(prog, 1e9)
	} else {
		fmt.Fprintln(w, "\n\n#ip", state.Arch.IP.Index)

		seen := make(map[string]struct{})

		for {
			fmt.Fprintln(w, state, prog[*state.IP])
			if !state.Step(prog) {
				break
			}

			if _, ok := seen[state.String()]; ok {
				log.Fatal("infinite loop")
			} else {
				seen[state.String()] = struct{}{}
			}
		}
	}

	return !state.Step(prog)
}

func puzzleInput19() (ipreg int, prog []wristdev.Instruction) {
	lines := PuzzleInputLines(19)

	var ipspec string
	ipspec, lines = lines[0], lines[1:]

	_, err := fmt.Sscanf(ipspec, "#ip %d", &ipreg)
	if err != nil {
		log.Fatal(err)
	}

	prog, err = wristdev.ParseProgram(lines)
	if err != nil {
		log.Fatal(err)
	}

	return ipreg, prog
}

/*

Program analysis:

  0  add ip 16 ip
  1  set 1 r4
  2  set 1 r5
  3  mul r4 r5 r1
  4  eq r1 r3 r1
  5  add r1 ip ip
  6  add ip 1 ip
  7  add r4 r0 r0
  8  add r5 1 r5
  9  gt r5 r3 r1
 10  add ip r1 ip
 11  set 2 ip
 12  add r4 1 r4
 13  gt r4 r3 r1
 14  add r1 ip ip
 15  set 1 ip
 16  mul ip ip ip
 17  add r3 2 r3
 18  mul r3 r3 r3
 19  mul ip r3 r3
 20  mul r3 11 r3
 21  add r1 6 r1
 22  mul r1 ip r1
 23  add r1 6 r1
 24  add r3 r1 r3
 25  add ip r0 ip
 26  set 0 ip
 27  set ip r1
 28  mul r1 ip r1
 29  add ip r1 r1
 30  mul ip r1 r1
 31  mul r1 14 r1
 32  mul r1 ip r1
 33  add r3 r1 r3
 34  set 0 r0
 35  set 0 ip

Annotation:

      ; goto L17
	  ;
  0         add ip 16 ip    ; jmp L17
      ; L1:
      ; r4 = 1
  1  L1:    set 1 r4
      ; L2:
      ; r5 = 1
  2  L2:    set 1 r5
      ; L3:
      ; if r4*r5 == r3 {
      ;   r0 += r4
      ; }  // r1 unused thereafter
  3  L3:    mul r4 r5 r1
  4         eq r1 r3 r1
  5         add r1 ip ip    ;         skip next if r1 == r3
  6         add ip 1 ip     ; skip next
  7         add r4 r0 r0    ; if r1 == r3
      ; r5++
  8         add r5 1 r5
      ; if r5 <= r3 {
      ;   goto L3
      ; }
  9         gt r5 r3 r1
 10         add ip r1 ip    ;        skip next if r5 > r3
 11         set 2 ip        ; jmp L3
      ; r4++
 12         add r4 1 r4
      ; if r4 <= r3 {
      ;   goto L2
      ; }
 13         gt r4 r3 r1
 14         add r1 ip ip    ;        skip next if r4 > r3
 15         set 1 ip        ; jmp L2
      ; return
	  ;
 16         mul ip ip ip
      ; L17:
      ; r3 = 2 * 2 * 19 * 11
 17  L17:   add r3 2 r3 ; line 17 reached after start only, so r3 == 2
 18         mul r3 r3 r3
 19         mul ip r3 r3 ; ip == 19
 20         mul r3 11 r3
      ; r1 = 6 * 22 + 6
 21         add r1 6 r1 ; line 17 reached after start only, so r1 == 6
 22         mul r1 ip r1 ; ip == 22
 23         add r1 6 r1
      ; r3 += r1
 24         add r3 r1 r3
      ; if r0 == 0 {
      ;   goto L1
      ; }
 25         add ip r0 ip
 26         set 0 ip
      ; r1 = (27*28 + 29) * 30 * 14 * 32
 27         set ip r1    ; ip == 27
 28         mul r1 ip r1 ; ip == 28
 29         add ip r1 r1 ; ip == 29
 30         mul ip r1 r1 ; ip == 30
 31         mul r1 14 r1 ; ip == 14
 32         mul r1 ip r1 ; ip == 32
      ; r3 += r1
 33         add r3 r1 r3
      ; r0 = 0
 34         set 0 r0
      ; goto L1
	  ;
 35         set 0 ip

Rewrite:
*/

func aoc19init(r0 int) (r3 int) {
	// lines 17..35
	r3 = 2 * 2 * 19 * 11
	r1 := 6*22 + 6
	r3 += r1
	if r0 != 0 {
		r1 = (27*28 + 29) * 30 * 14 * 32
		r3 += r1
	}

	return r3
}

func aoc19sim(r0 int) (res, in int) {

	r3 := aoc19init(r0)

	// loop for lines 1..15
	for r4 := 1; r4 <= r3; r4++ {

		// loop for lines 2..11
		for r5 := 1; r5 <= r3; r5++ {
			if r4*r5 == r3 {
				r0 += r4
			}
		}
	}

	return r0, r3
}

func aoc19divSum(v int) int {
	sum := 0
	for i := 1; i <= v; i++ {
		if v%i == 0 {
			sum += i
		}
	}
	return sum
	//fmt.Printf("%v: sum of divisors is %v\n", v, sum)
}
