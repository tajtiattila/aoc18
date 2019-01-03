package main

func DecomposePolymer(polymer string) string {
	return decomposePolymer([]byte(polymer))
}

func decomposePolymer(p []byte) string {
	q := make([]byte, 0, len(p))

	opp := func(u byte) byte {
		return u ^ 0x20 // switch case
	}

	for {
		final := true
		var lo byte // opposite of last unit
		for _, unit := range p {
			if unit == lo {
				final = false
				n := len(q) - 1
				q = q[:n]
				if n == 0 {
					lo = 0
				} else {
					lo = opp(q[n-1])
				}
			} else {
				q = append(q, unit)
				lo = opp(unit)
			}
		}

		if final {
			return string(q)
		}

		p, q = q, p
		q = q[:0]
	}

	panic("unreachable")
}

func CleanDecomposePolymer(polymer string) string {
	depol := func(u byte) byte {
		return u & ^byte(0x20)
	}
	units := make(map[byte]struct{})
	for i := 0; i < len(polymer); i++ {
		u := polymer[i]
		units[depol(u)] = struct{}{}
	}

	var best string
	for remove := range units {
		p := make([]byte, 0, len(polymer))
		for i := 0; i < len(polymer); i++ {
			u := polymer[i]
			if depol(u) != remove {
				p = append(p, u)
			}
		}

		dec := decomposePolymer(p)
		if best == "" || len(dec) < len(best) {
			best = dec
		}
	}
	return best
}
