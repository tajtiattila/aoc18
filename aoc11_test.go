package main

import "testing"

func TestFuelCellPower(t *testing.T) {
	tests := []struct {
		x, y, serial int

		n int

		want int
	}{
		// 1
		{122, 79, 57, 1, -5},
		{217, 196, 39, 1, 0},
		{101, 153, 71, 1, 4},
		{33, 45, 18, 3, 29},
		{21, 61, 42, 3, 30},
		// 2
		{90, 269, 18, 16, 113},
		{232, 251, 42, 12, 119},
	}

	for i, tt := range tests {
		got := FuelCellPowerN(tt.x, tt.y, tt.serial, tt.n)
		if got != tt.want {
			t.Errorf("%d: FuelCellPower(%d, %d, %d, %d) got %d, want %d",
				i, tt.x, tt.y, tt.serial, tt.n, got, tt.want)
		}
	}
}

func TestAoC11_1(t *testing.T) {
	tests := []struct {
		serial int

		x, y, power int
	}{
		{18, 33, 45, 29},
		{42, 21, 61, 30},
	}

	for i, tt := range tests {
		x, y, power := FindMaxFuelCellPowerN(300, tt.serial, 3)
		if tt.x != x || tt.y != y || tt.power != power {
			t.Errorf("%d: findmax1 %d got %d, %d, %d; want %d, %d, %d",
				i, tt.serial, x, y, power, tt.x, tt.y, tt.power)
		}
	}

	x, y, _ := FindMaxFuelCellPowerN(300, 1955, 3)
	t.Logf("1: max power for 1955 at: %d,%d", x, y)
}

func TestAoC11_2(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	tests := []struct {
		serial int

		x, y, n, power int
	}{
		{18, 90, 269, 16, 113},
		{42, 232, 251, 12, 119},
	}

	for i, tt := range tests {
		x, y, n, power := FindMaxFuelCellPowerAny(300, tt.serial)
		if tt.x != x || tt.y != y || tt.n != n || tt.power != power {
			t.Errorf("%d: findmax2 %d got %d, %d, %d, %d; want %d, %d, %d, %d",
				i, tt.serial, x, y, n, power, tt.x, tt.y, tt.n, tt.power)
		}
	}

	x, y, n, _ := FindMaxFuelCellPowerAny(300, 1955)
	t.Logf("2: max power for 1955 at: %d,%d,%d", x, y, n)
}
