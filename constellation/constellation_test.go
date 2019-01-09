package constellation

import (
	"fmt"
	"strings"
	"testing"
)

func TestConstellations(t *testing.T) {
	type test struct {
		wantc int
		src   string
	}

	tests := []test{
		{
			wantc: 2,
			src: `0,0,0,0
 3,0,0,0
 0,3,0,0
 0,0,3,0
 0,0,0,3
 0,0,0,6
 9,0,0,0
12,0,0,0`,
		},
		{
			wantc: 2,
			src: `0,0,0,0
 3,0,0,0
 0,3,0,0
 0,0,3,0
 0,0,0,3
 0,0,0,6
 9,0,0,0
12,0,0,0
 9,0,0,0`,
		},
		{
			wantc: 4,
			src: `-1,2,2,0
0,0,2,-2
0,0,0,-2
-1,2,0,0
-2,-2,-2,2
3,0,2,-1
-1,3,2,2
-1,0,-1,0
0,2,1,-2
3,0,0,0`,
		},
		{
			wantc: 3,
			src: `1,-1,0,1
2,0,-1,0
3,2,-1,0
0,0,3,1
0,0,-1,-1
2,3,-2,0
-2,2,0,0
2,-2,0,-1
1,-1,0,-1
3,2,0,2`,
		},
		{
			wantc: 8,
			src: `1,-1,-1,-2
-2,-2,0,1
0,2,1,3
-2,3,-2,1
0,2,3,-2
-1,-1,1,-2
0,-2,-1,0
-2,2,3,-1
1,2,2,0
-1,-2,0,-2`,
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprint(i+1), func(t *testing.T) {
			vsrc := strings.Split(tt.src, "\n")
			p, err := ParsePoints(vsrc)
			if err != nil {
				t.Fatal(err)
			}

			c := Constellations(p, 3)
			gotc := len(c)
			if gotc != tt.wantc {
				t.Fatalf("constellation: got %d, want %d", gotc, tt.wantc)
			}
		})
	}
}
