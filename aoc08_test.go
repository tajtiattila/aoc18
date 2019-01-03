package main

import (
	"strconv"
	"strings"
	"testing"
)

func TestAoC08(t *testing.T) {
	sample := parse08(t, "2 3 0 3 10 11 12 1 1 0 1 99 2 1 1 2")

	got1 := sample.SumMeta()
	want1 := 138
	if got1 != want1 {
		t.Fatalf("sample SumMeta got %v, want %v", got1, want1)
	}

	got2 := sample.Value()
	want2 := 66
	if got2 != want2 {
		t.Fatalf("sample Value got %v, want %v", got2, want2)
	}

	tree := parse08(t, input08)
	t.Logf("SumMeta: %d", tree.SumMeta())
	t.Logf("Value: %d", tree.Value())
}

func parse08(t *testing.T, src string) *Tree8 {
	var vi []int
	for _, s := range strings.Fields(src) {
		n, err := strconv.Atoi(s)
		if err != nil {
			t.Fatal(err)
		}
		vi = append(vi, n)
	}
	tree, err := ParseTree8(vi)
	if err != nil {
		t.Fatal(err)
	}
	return tree
}
