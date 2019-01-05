package resrsrch

import (
	"os"
	"strings"
	"testing"
)

func TestFlow(t *testing.T) {
	src := `
x=495, y=2..7
y=7, x=495..501
x=501, y=3..7
x=498, y=2..4
x=506, y=1..2
x=498, y=10..13
x=504, y=10..13
y=13, x=498..504`

	gs, err := ParseGroundSlice(strings.Split(src, "\n"))
	if err != nil {
		t.Fatal(err)
	}

	got := gs.Flood(500, 0, os.Stdout)

	// phase 1: 57, phase 2: 29
	want := FloodStat{Static: 29, Flow: 57 - 29}

	if got != want {
		t.Fatalf("got %v; want %v", got, want)
	}
}
