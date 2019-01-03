package main

import "github.com/pkg/errors"

type Tree8 struct {
	Child []*Tree8

	Meta []int
}

func ParseTree8(src []int) (*Tree8, error) {
	t, n, err := parseTree8(src)
	if err == nil && n != len(src) {
		err = errors.New("garbage after input")
	}
	if err != nil {
		err = errors.Wrapf(err, "Tree8 (at %d)", n)
	}
	return t, err
}

func parseTree8(src []int) (*Tree8, int, error) {
	if len(src) < 2 {
		return nil, 0, errors.New("short header")
	}
	t := &Tree8{}
	nc, nm := src[0], src[1]
	ofs := 2
	for i := 0; i < nc; i++ {
		child, n, err := parseTree8(src[ofs:])
		if err != nil {
			return nil, ofs + n, err
		}
		t.Child = append(t.Child, child)
		ofs += n
	}
	if len(src[ofs:]) < nm {
		return nil, ofs, errors.New("truncated meta")
	}
	t.Meta = src[ofs : ofs+nm]
	ofs += nm
	return t, ofs, nil
}

func (t *Tree8) SumMeta() int {
	sum := 0
	for _, child := range t.Child {
		sum += child.SumMeta()
	}
	for _, m := range t.Meta {
		sum += m
	}
	return sum
}

func (t *Tree8) Value() int {
	sum := 0

	if len(t.Child) == 0 {
		for _, m := range t.Meta {
			sum += m
		}
	} else {
		for _, m := range t.Meta {
			n := m - 1
			if n >= 0 && n < len(t.Child) {
				sum += t.Child[n].Value()
			}
		}
	}

	return sum
}
