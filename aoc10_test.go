package main

import (
	"fmt"
	"testing"
)

type skyray struct {
	x, y   int
	vx, vy int
}

func TestAoc10(t *testing.T) {
	var rays []skyray
	for _, s := range input10v {
		var r skyray
		//position=< 10703,  41994> velocity=<-1, -4>
		_, err := fmt.Sscanf(s, "position=<%d,%d> velocity=<%d,%d>",
			&r.x, &r.y, &r.vx, &r.vy)
		if err != nil {
			t.Fatalf("parse %q: %v", s, err)
		}
		rays = append(rays, r)
	}

	save := make([]skyray, len(rays))
	copy(save, rays)

	var niter, shown int
	adv := func(n int) {
		for i := range rays {
			r := &rays[i]
			r.x += n * r.vx
			r.y += n * r.vy
		}
		niter += n
	}

	lastscore := skyShapeRad(rays)
	for niter <= 1e6 {
		adv(1)
		score := skyShapeRad(rays)
		if score > lastscore {
			adv(-1)
			t.Logf("%6d %6d", niter, score)
			showSkyShape(t, rays)
			return
			showSkyShape(t, rays)
			shown++
			if shown > 1e2 {
				break
			}
		}
		lastscore = score
	}
}

func skyShapeRad(rays []skyray) int64 {
	var cx, cy int64
	for _, r := range rays {
		cx += int64(r.x)
		cy += int64(r.y)
	}
	cx /= int64(len(rays))
	cy /= int64(len(rays))

	var maxr2 int64
	for _, r := range rays {
		dx := int64(r.x) - cx
		dy := int64(r.y) - cy
		r2 := dx*dx + dy*dy
		if r2 > maxr2 {
			maxr2 = r2
		}
	}
	return maxr2
}

func showSkyShape(t *testing.T, rays []skyray) {
	ix := rays[0].x
	iy := rays[0].y
	ax, ay := ix, iy
	for _, r := range rays {
		if r.x < ix {
			ix = r.x
		}
		if r.y < iy {
			iy = r.y
		}
		if r.x > ax {
			ax = r.x
		}
		if r.y > ay {
			ay = r.y
		}
	}
	dx := ax - ix + 1
	dy := ay - iy + 1

	buf := make([]byte, dx*dy)
	for i := range buf {
		buf[i] = '.'
	}
	for _, r := range rays {
		x := r.x - ix
		y := r.y - iy
		ofs := x + y*dx
		buf[ofs] = '#'
	}

	for y := 0; y < dy; y++ {
		ofs := y * dx
		t.Logf("%s", buf[ofs:ofs+dx])
	}
}
