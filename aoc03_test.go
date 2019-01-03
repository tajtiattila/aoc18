package main

import (
	"fmt"
	"testing"
)

func TestAOC03_1(t *testing.T) {
	data := input03(t)
	t.Log(FindCutSpecOverlap(data, 2))
}

func TestAOC03_2(t *testing.T) {
	data := input03(t)
	t.Log(FindCutSpecSingleID(data))
}

func input03(t *testing.T) []CutSpec {
	v := make([]CutSpec, 0, len(input03v))
	for _, s := range input03v {
		var cs CutSpec
		_, err := fmt.Sscanf(s, "#%d @ %d,%d: %dx%d",
			&cs.ID, &cs.X, &cs.Y, &cs.Dx, &cs.Dy)
		if err != nil {
			t.Fatalf("scan %q: %v", s, err)
		}
		v = append(v, cs)
	}
	return v
}
