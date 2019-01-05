package pathfind

// Place represents a place in the problem space,
// such as a grid position or graph node pointer.
//
// A Place must be equality comparable.
type Place interface{}

type Space struct {
	// Adjacent finds places adjacent to p,
	// appending them to dst.
	Adjacent func(p Place, dst []Place) (adjacents []Place)

	// Step is called from Flood when p is reached
	// the first time.
	//
	// It reports if Flood should continue processing.
	//
	// When Step returns false the first time,
	// processing stops after processing the current distance.
	//
	// Step may be nil.
	Step func(p Place) (cont bool)
}

// Place -> distance
type FloodResult map[Place]int

func (m FloodResult) Distance(to Place) int {
	if d, ok := m[to]; ok {
		return d
	}
	return NotReached
}

const NotReached = int(^uint(0) >> 1)

func Flood(from Place, space Space) (res FloodResult, maxDist int) {

	if space.Step == nil {
		space.Step = func(p Place) bool { return true }
	}

	m := make(FloodResult)

	dist := 0
	active := []Place{from}
	m[from] = dist

	var nextactive []Place

	var adjbuf [16]Place
	adj := adjbuf[:]

	done := false
	for !done && len(active) != 0 {
		dist++
		for _, q := range active {
			adj = space.Adjacent(q, adj[:0])
			for _, p := range adj {
				if _, seen := m[p]; seen {
					continue
				}

				cont := space.Step(p)
				m[p] = dist
				nextactive = append(nextactive, p)
				maxDist = dist

				if !cont {
					done = true
				}
			}
		}

		active, nextactive = nextactive, active[:0]
	}

	return m, maxDist
}
