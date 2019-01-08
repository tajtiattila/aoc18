package nanobot

import (
	"fmt"
	"testing"
)

func TestMTree(t *testing.T) {
	type test struct {
		npoints int
		minp    point
		boxes   []MBox
	}

	tests := []test{
		{
			minp: point{x: 12, y: 12, z: 12},
			boxes: []MBox{
				Equidist(10, 12, 12, 2),
				Equidist(12, 14, 12, 2),
				Equidist(16, 12, 12, 4),
				Equidist(14, 14, 14, 6),
				Equidist(50, 50, 50, 200),
				Equidist(10, 10, 10, 5),
			},
		},
	}

	for tn, tt := range tests {
		t.Run(fmt.Sprintf("%d", tn+1), func(t *testing.T) {
			var bounds MBox
			for _, b := range tt.boxes {
				bounds = bounds.Extend(b)
			}

			tree := MTree{Bounds: bounds}
			for _, b := range tt.boxes {
				tree.Add(b)
			}

			var best *MTree
			tree.WalkLeaves(func(node *MTree) {
				if best == nil || node.Count > best.Count {
					best = node
				}
			})

			x, y, z, ok := best.Bounds.MinPoint()
			gotp := point{x: x, y: y, z: z}
			if !ok || gotp != tt.minp {
				t.Fatalf("got %v min point %v; want %v", ok, gotp, tt.minp)
			}
		})
	}
}
