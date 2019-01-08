package nanobot

import "math/bits"

type Bot struct {
	X, Y, Z int
	Radius  int
}

type word uint64

const (
	MaxBot = nword * wordBits

	wantBot  = 1000
	wordBits = 64 // bitlen(word); power of 2
	nword    = (wantBot + wordBits - 1) / wordBits

	maskBit  = wordBits - 1
	maskWord = int(^0) ^ maskBit
)

type Botset struct {
	bits  [nword]word
	count int
}

func (s *Botset) Set(n int) {
	i, m := n/wordBits, word(1)<<uint(n&maskBit)
	if s.bits[i]&m == 0 {
		s.bits[i] |= m
		s.count++
	}
}

func (s *Botset) Get(n int) bool {
	i, m := n/wordBits, word(1)<<uint(n%maskBit)
	return s.bits[i]&m != 0
}

// first bit set; MaxBot if none
func (s *Botset) First() int {
	return s.next(0)
}

// first bit set after n; MaxBot if none
func (s *Botset) Next(n int) int {
	return s.next(n + 1)
}

func (s *Botset) next(n int) int {
	i, o := n/wordBits, n&maskWord
	m := ^word(0) << uint(n%maskBit)
	for i < nword {
		z := bits.TrailingZeros64(uint64(s.bits[i] & m))
		if z < wordBits {
			return o + z
		}
		i++
		o += wordBits
		m = ^word(0)
	}
	return MaxBot
}
