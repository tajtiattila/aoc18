package gridregexp

import (
	"bytes"
	"fmt"
	"testing"
)

func TestParse(t *testing.T) {
	type test struct {
		src      string
		ok       bool
		maxDoors int
	}
	tests := []test{
		test{
			src: "^ENWWW(NEEE|SSE(EE|N))$",
			ok:  true,
		},
		test{
			src: "^ENNWSWW(NEWS|)SSSEEN(WNSE|)EE(SWEN|)NNN$",
			ok:  true,
		},
		test{
			src:      "^ESSWWN(E|NNENN(EESS(WNSE|)SSS|WWWSSSSE(SW|NNNE)))$",
			ok:       true,
			maxDoors: 23,
		},
		test{
			src:      "^WSSEESWWWNW(S|NENNEEEENN(ESSSSW(NWSW|SSEN)|WSWWN(E|WWS(E|SS))))$",
			ok:       true,
			maxDoors: 31,
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			expr, err := Parse(tt.src)
			if tt.ok {
				if err != nil {
					t.Fatal("parse failed:", err)
				}
				if expr == nil {
					t.Fatal("parse expr is nil")
				}

				sexpr := expr.String()
				if sexpr != tt.src {
					t.Fatalf("expr is %q, want %q", sexpr, tt.src)
				}

				var buf bytes.Buffer
				m := expr.Map()
				if err := m.Write(&buf); err != nil {
					t.Fatal("write map", err)
				}
				t.Logf("%q:\n%s", tt.src, buf.String())

				md := m.MaxDoors()
				if tt.maxDoors != 0 && tt.maxDoors != md {
					t.Fatalf("got max doors %v, want %v", md, tt.maxDoors)
				}
			} else {
				if err != nil {
					t.Fatal("parse successful but should have failed")
				}
			}
		})
	}
}
