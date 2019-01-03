package main

import (
	"strings"
	"testing"
)

func TestAoC7(t *testing.T) {
	sample := parseAssemblyLinks(t, `Step C must be finished before step A can begin.
Step C must be finished before step F can begin.
Step A must be finished before step B can begin.
Step A must be finished before step D can begin.
Step B must be finished before step E can begin.
Step D must be finished before step E can begin.
Step F must be finished before step E can begin.`)

	got1 := strings.Join(SortAssemblyInstr(sample), "")
	want1 := "CABDFE"
	if got1 != want1 {
		t.Logf("SortAssemblyInstr: got %q, want %q", got1, want1)
	}

	got2 := TimeAssembly(sample, 2, func(s string) int {
		if len(s) != 1 {
			t.Fatal("timeinstrfunc logic error")
		}
		return int(s[0] - 'A' + 1)
	})
	want2 := 15
	if got2 != want2 {
		t.Logf("TimeAssembly: got %v, want %v", got2, want2)
	}

	data := parseAssemblyLinks(t, input07)
	t.Log("SortAssemblyInstr:", strings.Join(SortAssemblyInstr(data), ""))
	t.Log("TimeAssembly:", TimeAssembly(data, 5, func(s string) int {
		if len(s) != 1 {
			t.Fatal("invalid data")
		}
		return int(s[0] - 'A' + 61)
	}))
}

func parseAssemblyLinks(t *testing.T, src string) []AssemblyLink {
	var v []AssemblyLink
	for _, l := range strings.Split(src, "\n") {
		al, err := ParseAssemblyLink(l)
		if err != nil {
			t.Fatal(err)
		}
		v = append(v, al)
	}
	return v
}
