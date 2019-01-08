package nanobot

import (
	"fmt"
	"sort"
)

const maxsub = 16

type MTree struct {
	Bounds MBox

	Sub []MBox

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

		t.Sub = append(t.Sub, box)
		if len(t.Sub) < maxsub {
			return
		}

		t.split()
	}

	for i := range t.Child {
		c := &t.Child[i]
		c.Add(box)
	}
}

func (t *MTree) split() {
	splits := make([][]int, 4)
	nchild := 1
	for axis := range splits {
		var v []int
		v = append(v, t.Bounds.Min[axis], t.Bounds.Max[axis])

		for _, sub := range t.Sub {
			if t.Bounds.Min[axis] < sub.Min[axis] {
				v = append(v, sub.Min[axis])
			}
			if sub.Max[axis] < t.Bounds.Max[axis] {
				v = append(v, sub.Max[axis])
			}
		}
		sort.Ints(v)

		// dedup
		j, lastc := 1, v[0]
		for _, c := range v[1:] {
			if c != lastc {
				v[j] = c
				j++
				lastc = c
			}
		}

		v = v[:j]
		splits[axis] = v
		nchild *= len(v) - 1
	}
	fmt.Printf("%p split %v\n", t, nchild)

	if nchild < 2 {
		panic("splitting empty node")
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

	for i := range t.Child {
		c := &t.Child[i]
		for _, sub := range t.Sub {
			x := c.Bounds.Intersect(sub)
			if x != c.Bounds && !x.Empty() {
				panic("invalid split")
			}
			if !x.Empty() {
				c.Count++
			}
		}
	}
	t.Sub = nil
}

func (t *MTree) WalkLeaves(f func(node *MTree)) {
	if t.Child == nil {
		if len(t.Sub) <= 1 {
			f(t)
			return
		}

		t.split()
	}

	for i := range t.Child {
		c := &t.Child[i]
		c.WalkLeaves(f)
	}
}
