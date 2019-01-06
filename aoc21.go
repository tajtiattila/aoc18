package main

import (
	"context"
	"fmt"
	"log"

	"github.com/tajtiattila/aoc18/wristdev"
)

func wristdev21() {
	ipreg, prog := programInput(21)

	const nreg = 6

	regnames := make([]string, nreg)
	for i := range regnames {
		regnames[i] = fmt.Sprintf("r%d", i)
	}
	regnames[ipreg] = "ip"

	if false {
		fmt.Println("disassembly:")
		for i, l := range prog {
			fmt.Printf("%3d  %s\n", i, l.Fmt(regnames))
		}
		fmt.Println()
	}

	// verify simluation by comparing output with that of the program
	fns := []aoc21simfunc{
		aoc21simprog,
		aoc21sim0,
		aoc21sim1,
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	const (
		simverifysteps = 1e6
	)

	ch := merge21(ctx, 0, fns)
	step := 0
	for ss := range ch {
		step++
		prog := ss[0]
		for i := 1; i < len(ss); i++ {
			if ss[i] != prog {
				log.Printf("step %d: state #%d differ\n", step, i)
				log.Fatal(ss)
			}
		}
		if step == simverifysteps {
			return
		}
	}

	const (
		c0 = 14906355 // 0xe373f3
		c1 = 65899

		m24 = 0xffffff
	)

	if verbose {
		div := 0
		for i := 2; i*i <= c1; i++ {
			if c1%i == 0 {
				div = i
				break
			}
		}
		if div != 0 {
			fmt.Printf("%v is divisible by %v\n", c1, div)
		} else {
			fmt.Printf("%v is prime\n", c1)
		}
	}

	r3m := make(map[uint32]struct{})
	r3 := uint32(0)
	var lastr3 uint32
	for step := 0; ; step++ {
		r1 := uint32(r3) | 0x10000
		r3 = uint32(c0)

		for r1 != 0 { // outer loop
			r3 = (r3 + (r1 & 0xff)) & 0xffffff
			r3 = (r3 * c1) & 0xffffff
			r1 /= 256 // inner loop
		}

		if step == 0 {
			fmt.Println("21/1:", r3)
		}

		if _, seen := r3m[r3]; !seen {
			r3m[r3] = struct{}{}
			lastr3 = r3
		} else {
			fmt.Println("21/2:", lastr3)
			break
		}
	}
}

type aoc21simfunc func(ctx context.Context, r0 regt) <-chan simstate

func merge21(ctx context.Context, r0 regt, fns []aoc21simfunc) <-chan []simstate {

	chans := make([]<-chan simstate, len(fns))
	for i, fn := range fns {
		chans[i] = fn(ctx, r0)
	}

	ch := make(chan []simstate)
	go func() {
		defer close(ch)

		r := make([]simstate, len(chans))

		for i := range r {
			if chans[i] == nil {
				continue
			}

			select {
			case <-ctx.Done():
				return
			case s, ok := <-chans[i]:
				r[i] = s
				if !ok {
					chans[i] = nil
				}
			}
		}

		select {
		case <-ctx.Done():
			return
		case ch <- r:
		}

	}()

	return ch
}

func aoc21simprog(ctx context.Context, r0 regt) <-chan simstate {
	ipreg, prog := programInput(21)

	const nreg = 6

	arch := wristdev.ArchWithIP(nreg, ipreg)
	state := arch.State(int(r0))

	ch := make(chan simstate)
	go func() {
		defer close(ch)

		for {
			if !state.Step(prog) {
				return
			}
			if *state.IP == 28 {
				var st simstate
				for i := 0; i < len(st); i++ {
					st[i] = regt(state.R[i])
				}
				select {
				case <-ctx.Done():
					return
				case ch <- st:
				}
			}
		}
	}()

	return ch
}

/*

Disassembly:
  0  set 123 r3
  1  ban r3 456 r3
  2  eq r3 72 r3
  3  add r3 ip ip
  4  set 0 ip
  5  set 0 r3
  6  bor r3 65536 r1
  7  set 14906355 r3
  8  ban r1 255 r4
  9  add r3 r4 r3
 10  ban r3 16777215 r3
 11  mul r3 65899 r3
 12  ban r3 16777215 r3
 13  gt 256 r1 r4
 14  add r4 ip ip
 15  add ip 1 ip
 16  set 27 ip
 17  set 0 r4
 18  add r4 1 r2
 19  mul r2 256 r2
 20  gt r2 r1 r2
 21  add r2 ip ip
 22  add ip 1 ip
 23  set 25 ip
 24  add r4 1 r4
 25  set 17 ip
 26  set r4 r1
 27  set 7 ip
 28  eq r3 r0 r4
 29  add r4 ip ip
 30  set 5 ip

Analysis:

	; // check 126&456 == 72
	; // uses r3 only, then sets it to 0, therefore no-op
  0         set 123 r3
  1         ban r3 456 r3
  2         eq r3 72 r3
  3         add r3 ip ip
  4         set 0 ip
  5         set 0 r3
    ; Line6: // begin main loop
	; r1 = r3 | 0x10000
  6         bor r3 65536 r1 ; 65536 == 0x10000
	; r3 = 14906355
  7         set 14906355 r3 ; 14906355 == 0xe373f3
	; Line8: // begin outer loop
	; r4 = r1 & 0xff
	; r3 = (r3 + r4) & 0xffffff
	; r3 = (r3 * 65899) & 0xffffff
  8         ban r1 255 r4 ; 255 == 0xff
  9         add r3 r4 r3
 10         ban r3 16777215 r3 ; 16777215 == 0xffffff
 11         mul r3 65899 r3
 12         ban r3 16777215 r3 ; 16777215 == 0xffffff
	; if r1 < 256 {
	;   goto Line28
	; }
 13         gt 256 r1 r4
 14         add r4 ip ip
 15         add ip 1 ip
 16         set 27 ip
	; r4 = 0
 17         set 0 r4
	; Line18: // begin inner loop
	; r2 = (r4 + 1) * 256
 18         add r4 1 r2
 19         mul r2 256 r2
	; if r2 > r1 {
	;   r2 = 1
	;   goto Line26
	; }
 20         gt r2 r1 r2
 21         add r2 ip ip
 22         add ip 1 ip
 23         set 25 ip
	; r4++
 24         add r4 1 r4
	; goto Line18 // end inner loop
 25         set 17 ip
	; Line26:
	; r1 = r4
 26         set r4 r1
	; goto Line8 // end outer loop
 27         set 7 ip
	; Line28:
	; if r3 != r0 {
	;   goto Line6
	; } // end main loop
 28         eq r3 r0 r4
 29         add r4 ip ip
 30         set 5 ip

*/

type regt int

type simstate [5]regt

func yieldsimstate(ctx context.Context, ch chan simstate,
	r0, r1, r2, r3, r4 regt) bool {
	st := simstate{r0, r1, r2, r3, r4}
	select {
	case <-ctx.Done():
		return true
	case ch <- st:
		return false
	}
}

func aoc21sim0(ctx context.Context, r0 regt) <-chan simstate {

	ch := make(chan simstate)
	go func() {
		defer close(ch)

		var r1, r2, r3, r4 regt

	Line6: // begin main loop
		r1 = r3 | 0x10000
		r3 = 14906355

	Line8: // begin outer loop
		r4 = r1 & 0xff
		r3 = (r3 + r4) & 0xffffff
		r3 = (r3 * 65899) & 0xffffff
		if r1 < 256 {
			goto Line28
		}
		r4 = 0

	Line18: // begin inner loop
		r2 = (r4 + 1) * 256
		if r2 > r1 {
			r2 = 1
			goto Line26
		}
		r4++
		goto Line18 // end inner loop

	Line26:
		r1 = r4
		goto Line8 // end outer loop

	Line28:
		if yieldsimstate(ctx, ch, r0, r1, r2, r3, r4) {
			return
		}
		if r3 != r0 {
			goto Line6
		} // end main loop
	}()

	return ch
}

func aoc21sim1(ctx context.Context, r0 regt) <-chan simstate {

	ch := make(chan simstate)
	go func() {
		defer close(ch)

		var r1, r2, r3, r4 regt

		// check 126&456 == 72

		// r2 is not used below, and is always 1 on line 28
		// r3 is set to ensure the loop starts
		r2, r3 = 1, r0|0x10000

		for r3 != r0 {

			r1 = r3 | 0x10000
			r3 = 14906355

			// simr1 needed to satisfy
			var simr1 regt

			for r1 != 0 {
				simr1 = r1
				r3 = (r3 + (r1 & 0xff)) & 0xffffff
				r3 = (r3 * 65899) & 0xffffff
				r1 /= 256 // inner loop
			}

			// satisy simstate
			r1 = simr1
			r4 = r1

			if yieldsimstate(ctx, ch, r0, r1, r2, r3, r4) {
				return
			}
		}
	}()

	return ch
}
