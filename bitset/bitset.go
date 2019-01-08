package bitset

import "math/bits"

type word uint

const (
	wordBits = bits.UintSize

	maskBit  = wordBits - 1
	maskWord = int(^0) ^ maskBit
)

type Bitset struct {
	p     []word // bit storage
	count int    // number of bits set
}

// Ones returns a new Bitset with bits [0..n) set.
func Ones(n int) Bitset {
	if n == 0 {
		return Bitset{}
	}
	s := Bitset{
		p:     make([]word, (n+wordBits-1)/wordBits),
		count: n,
	}
	bit, i := 0, 0
	for bit+wordBits <= n {
		s.p[i] = ^word(0)
		bit = bit + wordBits
		i++
	}
	if i < n {
		s.p[i] = ^word(0) >> (wordBits - uint(i))
	}
	return s
}

// Count returns the number of set bits
func (s Bitset) Count() int {
	return s.count
}

func (s Bitset) Clone() Bitset {
	c := Bitset{
		p:     s.p,
		count: s.count,
	}
	copy(c.p, s.p)
	return c
}

func (s *Bitset) Reset() {
	s.p = s.p[:0]
	s.count = 0
}

func (s *Bitset) haveWord(i int) {
	n := i + 1
	if len(s.p) > n {
		return
	}
	if cap(s.p) > n {
		i = len(s.p)
		s.p = s.p[:n]
		for ; i < n; i++ {
			s.p[i] = 0
		}
		return
	}

	const minlen = 16
	if n < minlen {
		n = minlen
	}

	if m := 2 * len(s.p); m > n {
		n = m
	}

	p := make([]word, n)
	copy(p, s.p)
	s.p = p
}

func (s *Bitset) Set(n int) {
	i, m := n/wordBits, word(1)<<uint(n&maskBit)
	s.haveWord(i)
	if s.p[i]&m == 0 {
		s.p[i] |= m
		s.count++
	}
}

func (s *Bitset) Clear(n int) {
	i, m := n/wordBits, word(1)<<uint(n&maskBit)
	if i >= len(s.p) {
		return
	}
	if s.p[i]&m != 0 {
		s.p[i] &= ^m
		s.count--
	}
}

func (s Bitset) Get(n int) bool {
	i, m := n/wordBits, word(1)<<uint(n&maskBit)
	return s.p[i]&m != 0
}

// first bit set; -1 if none
func (s Bitset) First() int {
	return s.next(0)
}

// first bit set after n; -1 if none
func (s Bitset) Next(n int) int {
	return s.next(n + 1)
}

func (s Bitset) next(n int) int {
	i, o := n/wordBits, n&maskWord
	m := ^word(0) << uint(n&maskBit)
	for i < len(s.p) {
		z := bits.TrailingZeros64(uint64(s.p[i] & m))
		if z < wordBits {
			return o + z
		}
		i++
		o += wordBits
		m = ^word(0)
	}
	return -1
}
