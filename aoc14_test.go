package main

import (
	"bytes"
	"testing"
)

func TestAoC14_1(t *testing.T) {
	cr := chocReciper{
		recipes: []byte{3, 7},
		elf1:    0,
		elf2:    1,
	}

	tests := []struct {
		after int
		want  string
	}{
		{5, "0124515891"},
		{9, "5158916779"},
		{18, "9251071085"},
		{2018, "5941429882"},
	}

	for _, tt := range tests {
		got := cr.after(tt.after, 10)
		if got != tt.want {
			t.Errorf("after %v got %v; want %v", tt.after, got, tt.want)
		}
	}

	const after = 430971
	t.Logf("after %v: %v", after, cr.after(after, 10))
}

func TestAoC14_2(t *testing.T) {
	cr := chocReciper{
		recipes: []byte{3, 7},
		elf1:    0,
		elf2:    1,
	}

	tests := []struct {
		find string
		want int
	}{
		{"51589", 9},
		{"01245", 5},
		{"92510", 18},
		{"59414", 2018},
	}

	for _, tt := range tests {
		got := cr.findIndex(tt.find)
		if got != tt.want {
			t.Errorf("find %v got %v; want %v", tt.find, got, tt.want)
		}
	}

	const find = "430971"
	t.Logf("find %v: %v", find, cr.findIndex(find))
}

type chocReciper struct {
	recipes    []byte
	elf1, elf2 int
}

func (r *chocReciper) step() {
	score1 := r.recipes[r.elf1]
	score2 := r.recipes[r.elf2]

	s := score1 + score2
	if s >= 10 {
		r.recipes = append(r.recipes, s/10)
	}
	r.recipes = append(r.recipes, s%10)

	r.elf1 = (r.elf1 + 1 + int(score1)) % len(r.recipes)
	r.elf2 = (r.elf2 + 1 + int(score2)) % len(r.recipes)
}

func (r *chocReciper) after(i, n int) string {
	for len(r.recipes) < i+n {
		r.step()
	}

	buf := make([]byte, n)
	for j := 0; j < n; j++ {
		buf[j] = r.recipes[i+j] + '0'
	}
	return string(buf)
}

func (r *chocReciper) findIndex(s string) int {
	b := make([]byte, len(s))
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c < '0' || c > '9' {
			return -1
		}
		b[i] = byte(c) - '0'
	}

	for len(r.recipes) < len(b) {
		r.step()
	}

	ofs := 0
	for len(r.recipes) < 1e8 {
		i := bytes.Index(r.recipes[ofs:], b)
		if i >= 0 {
			return ofs + i
		}

		ofs = len(r.recipes) - len(b)
		n := len(r.recipes) * 2
		for len(r.recipes) < n {
			r.step()
		}
	}
	return -1
}
