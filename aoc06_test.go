package main

import (
	"fmt"
	"testing"
)

func TestAoC06(t *testing.T) {
	sample := []Point{
		Pt(1, 1),
		Pt(1, 6),
		Pt(8, 3),
		Pt(3, 4),
		Pt(5, 5),
		Pt(8, 9),
	}

	got1 := FindNonInfArea(sample)
	want1 := 17
	if got1 != want1 {
		t.Fatalf("FindNonInfArea: got %d, want %d", got1, want1)
	}

	got2 := FindAreaCloserThan(sample, 32)
	want2 := 16
	if got2 != want2 {
		t.Fatalf("FindAreaCloserThan: got %d, want %d", got2, want2)
	}

	var data []Point
	for _, s := range input06v {
		var x, y int
		if _, err := fmt.Sscanf(s, "%d, %d", &x, &y); err != nil {
			t.Fatalf("parse %q: %v", s, err)
		}
		data = append(data, Pt(x, y))
	}

	t.Log(FindNonInfArea(data))
	t.Log(FindAreaCloserThan(data, 10000))
}
