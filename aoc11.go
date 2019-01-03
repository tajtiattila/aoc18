package main

func FuelCellPower(x, y, serial int) int {
	rackid := x + 10
	pow := rackid*y + serial
	pow *= rackid

	pl := (pow % 1000) / 100
	return pl - 5
}

func FuelCellPowerN(x, y, serial, n int) int {
	return fuelCellPowerXY(x, y, serial, n, n)
}

func fuelCellPowerXY(x, y, serial, nx, ny int) int {
	sum := 0
	for dx := 0; dx < nx; dx++ {
		for dy := 0; dy < ny; dy++ {
			sum += FuelCellPower(x+dx, y+dy, serial)
		}
	}
	return sum
}

func FindMaxFuelCellPowerN(dim, serial, n int) (left, top, power int) {
	x, y := 1, 1
	right := true // moving left
	p := FuelCellPowerN(x, y, serial, n)

	mx, my, mp := x, y, p
	for {
		if p > mp {
			mx, my, mp = x, y, p
		}

		movedown := false
		if right {
			if x+n < dim {
				// move window right
				p -= fuelCellPowerXY(x, y, serial, 1, n)
				p += fuelCellPowerXY(x+n, y, serial, 1, n)
				x++
			} else {
				movedown = true
			}
		} else {
			if x > 1 {
				// move window left
				x--
				p -= fuelCellPowerXY(x+n, y, serial, 1, n)
				p += fuelCellPowerXY(x, y, serial, 1, n)
			} else {
				movedown = true
			}
		}

		if movedown {
			if y+n >= dim {
				return mx, my, mp
			}

			// move down and swap left/right direction
			right = !right
			p -= fuelCellPowerXY(x, y, serial, n, 1)
			p += fuelCellPowerXY(x, y+n, serial, n, 1)
			y++
		}
	}

	panic("not reached")
}

func FindMaxFuelCellPowerAny(dim, serial int) (x, y, n, power int) {
	var mx, my, mn, mp int
	for n := 3; n <= dim; n++ {
		x, y, p := FindMaxFuelCellPowerN(dim, serial, n)
		if p > mp {
			mx, my, mn, mp = x, y, n, p
		}
	}
	return mx, my, mn, mp
}
