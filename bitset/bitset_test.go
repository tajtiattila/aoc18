package bitset

import (
	"fmt"
	"testing"
)

func TestBitset(t *testing.T) {
	tests := [][]int{
		[]int{0, 50, 128, 256},
		[]int{1, 2, 3, 4, 5, 6},
		[]int{63, 64},
		[]int{999},
		[]int{9999}, // because of haveWord.minlen
	}

	for i, tt := range tests {
		t.Run(fmt.Sprint(i), func(t *testing.T) {

			t.Log(tt)

			var bs Bitset
			for _, bit := range tt {
				bs.Set(bit)
			}

			i := 0
			for bit := bs.First(); bit > 0; bit = bs.Next(bit) {
				t.Log(i, bit)
				if tt[i] != bit {
					t.Fatalf("mismatch at %v: %v != %v", i, tt[i], bit)
				}
				i++
			}

			bs.Reset()
			for _, bit := range tt {
				bs.Set(bit)
			}

			i = 0
			for bit := bs.First(); bit > 0; bit = bs.Next(bit) {
				if tt[i] != bit {
					t.Fatalf("mismatch at %v: %v != %v", i, tt[i], bit)
				}
				i++
			}

			const ones = 100
			bs = Ones(100)
			for _, bit := range tt {
				bs.Clear(bit)
			}
			gotc := bs.Count()
			wantc := ones
			for _, bit := range tt {
				if bit < ones {
					wantc--
				}
			}
			if gotc != wantc {
				t.Fatalf("got count %v, want %v", gotc, wantc)
			}
		})
	}
}
