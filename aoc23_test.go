package main

import (
	"math/bits"
	"math/rand"
	"testing"
)

func benchbitvalues() []uint64 {
	s := rand.New(rand.NewSource(0))

	const n = 64

	v := make([]uint64, 0, 2*n)

	for i := 0; i < n; i++ {
		r := s.Uint64()
		v = append(v, r, bits.Reverse64(r))
	}

	return v
}

var benchbitsresult int

func BenchmarkLeadingZeros64(b *testing.B) {

	vv := benchbitvalues()

	n := 0

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, v := range vv {
			n += bits.LeadingZeros64(v)
		}
	}

	benchbitsresult = n
}

func BenchmarkTrailingZeros64(b *testing.B) {

	vv := benchbitvalues()

	n := 0

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, v := range vv {
			n += bits.TrailingZeros64(v)
		}
	}

	benchbitsresult = n
}
