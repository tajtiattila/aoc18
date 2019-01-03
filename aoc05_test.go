package main

import "testing"

func TestAoC05(t *testing.T) {
	sample := "dabAcCaCBAcCcaDA"
	got1 := DecomposePolymer(sample)
	want1 := "dabCBAcaDA"
	if got1 != want1 {
		t.Fatalf("sample %v decompose: got %v want %v", sample, got1, want1)
	}

	got2 := CleanDecomposePolymer(sample)
	want2 := "daDA"
	if got2 != want2 {
		t.Fatalf("sample %v clean-decompose: got %v want %v", sample, got2, want2)
	}

	d := DecomposePolymer(input05)
	t.Logf("decompose len(%d) -> len(%d)", len(input05), len(d))

	cd := CleanDecomposePolymer(input05)
	t.Logf("clean-decompose len(%d) -> len(%d)", len(input05), len(cd))
}
