package main

import "testing"

func TestAoC13Sample(t *testing.T) {
	sample := []string{
		`/->-\        `,
		`|   |  /----\`,
		`| /-+--+-\  |`,
		`| | |  | v  |`,
		`\-+-/  \-+--/`,
		`  \------/   `,
	}
	m, err := ParseMinecartMap(sample)
	if err != nil {
		t.Fatal(err)
	}

	for i := 0; i < 50; i++ {
		m.Tick()
	}

	got1 := m.crash[0]
	want1 := Pt(7, 3)
	if got1 != want1 {
		t.Errorf("1st crash is at %d,%d; want %d,%d", got1.X, got1.Y, want1.X, want1.Y)
	}
}

func TestAoC13Sample2(t *testing.T) {
	sample := []string{
		`/>-<\  `,
		`|   |  `,
		`| /<+-\`,
		`| | | v`,
		`\>+</ |`,
		`  |   ^`,
		`  \<->/`,
	}
	m, err := ParseMinecartMap(sample)
	if err != nil {
		t.Fatal(err)
	}

	const nticks = 10
	for i := 0; i < nticks; i++ {
		m.Tick()
		m.sortcarts()
		if len(m.cart) < 2 {
			break
		}
	}

	if len(m.cart) != 1 {
		t.Errorf("%d carts running after %d ticks", len(m.cart), nticks)
	}

	cart := m.cart[0]
	got1 := Pt(cart.x, cart.y)
	want1 := Pt(6, 4)
	if got1 != want1 {
		t.Errorf("last cart is at %d,%d; want %d,%d", got1.X, got1.Y, want1.X, want1.Y)
	}
}

func TestAoC13(t *testing.T) {
	m, err := ParseMinecartMap(input13v)
	if err != nil {
		t.Fatal(err)
	}

	showncrash := false
	for i := 0; ; i++ {
		m.Tick()

		if len(m.crash) > 0 && !showncrash {
			pt := m.crash[0]
			t.Logf("tick %d: first crash at %d,%d", i, pt.X, pt.Y)
			showncrash = true
		}

		if len(m.cart) == 1 {
			cart := m.cart[0]
			t.Logf("tick %d: last cart at %d,%d", i, cart.x, cart.y)
			return
		}
	}
}
