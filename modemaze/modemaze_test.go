package modemaze

import (
	"bytes"
	"testing"
)

func TestSample(t *testing.T) {
	m := New(510, 10, 10)

	var buf bytes.Buffer
	m.Write(&buf, 16, 16)
	t.Logf("\n%s\n", buf.String())

	got := m.RiskLevel()
	want := 114

	if got != want {
		t.Fatalf("got risk level %v; want %v", got, want)
	}

	got = m.PathDuration()
	want = 45

	if got != want {
		t.Fatalf("got duration %v; want %v", got, want)
	}
}
