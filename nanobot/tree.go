package nanobot

type MTree struct {
	Bounds MBox

	Child []MTree

	Count int
}

func (t *MTree) Add(box MBox) {
	x := t.Bounds.Intersect(box)
	if x.Empty() {
		return
	}

	if t.Child == nil {

		if x == t.Bounds {
			t.Count++
			return
		}

		splits := make([][]int, 4)
		nchild := 1
		for axis := range splits {
			var v []int
			v = append(v, t.Bounds.Min[axis])
			if t.Bounds.Min[axis] < x.Min[axis] {
				v = append(v, x.Min[axis])
			}
			if x.Max[axis] < t.Bounds.Max[axis] {
				v = append(v, x.Max[axis])
			}
			v = append(v, t.Bounds.Max[axis])

			splits[axis] = v
			nchild *= len(v) - 1
		}

		t.Child = make([]MTree, nchild)

		var cbox MBox
		idx := [4]int{1, 1, 1, 1}

		for i := range t.Child {
			for axis := range splits {
				cbox.Min[axis] = splits[axis][idx[axis]-1]
				cbox.Max[axis] = splits[axis][idx[axis]]
			}

			t.Child[i] = MTree{
				Bounds: cbox,
				Count:  t.Count,
			}

			for axis := range splits {
				idx[axis]++
				if idx[axis] < len(splits[axis]) {
					break
				} else {
					idx[axis] = 1
				}
			}
		}
	}

	for i := range t.Child {
		c := &t.Child[i]
		c.Add(box)
	}
}

func (t *MTree) WalkLeaves(f func(*MTree)) {
	if t.Child == nil {
		f(t)
		return
	}

	for i := range t.Child {
		c := &t.Child[i]
		c.WalkLeaves(f)
	}
}
