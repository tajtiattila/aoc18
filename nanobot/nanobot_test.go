package nanobot

import (
	"fmt"
	"testing"
)

func TestBotset(t *testing.T) {
	tests := [][]int{
		[]int{0, 50, 128, 256},
		[]int{1, 2, 3, 4, 5, 6},
		[]int{63, 64},
		[]int{999},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			var bs Botset
			for _, bit := range tt {
				bs.Set(bit)
			}

			t.Log(tt)

			i := 0
			for bit := bs.First(); bit < MaxBot; bit = bs.Next(bit) {
				t.Log(i, bit)
				if tt[i] != bit {
					t.Fatalf("mismatch at %v: %v != %v", i, tt[i], bit)
				}
				i++
			}
		})
	}
}
