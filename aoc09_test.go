package main

import "testing"

func TestAoC9(t *testing.T) {
	tests := []struct {
		nplayers, nmarbles int

		wantplayer, wantscore int
	}{
		{9, 25, 5, 32},
		{10, 1618, 0, 8317},
		{13, 7999, 0, 146373},
		{17, 1104, 0, 2764},
		{21, 6111, 0, 54718},
		{30, 5807, 0, 37305},
	}

	for _, x := range tests {
		gotplayer, gotscore := PlayMarbleRing(x.nplayers, x.nmarbles)

		if x.wantplayer != 0 && gotplayer != x.wantplayer {
			t.Errorf("%d players, %d marbles: got winner %d, want %d",
				x.nplayers, x.nmarbles, gotplayer, x.wantplayer)
		}

		if gotscore != x.wantscore {
			t.Errorf("%d players, %d marbles: got score %d, want %d",
				x.nplayers, x.nmarbles, gotscore, x.wantscore)
		}
	}

	// 452 players; last marble is worth 71250 points
	nplayers, nmarbles := 452, 71250
	winner, score := PlayMarbleRing(nplayers, nmarbles)
	t.Logf("PlayMarbleRing 1: %d players, %d marbles: player %d wins with %d score",
		nplayers, nmarbles, winner, score)

	nmarbles *= 100
	winner, score = PlayMarbleRing(nplayers, nmarbles)
	t.Logf("PlayMarbleRing 2: %d players, %d marbles: player %d wins with %d score",
		nplayers, nmarbles, winner, score)
}
