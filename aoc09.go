package main

type mring struct {
	left, right *mring

	num int
}

func PlayMarbleRing(nplayers, nmarbles int) (winner, score int) {
	cur := &mring{}
	cur.left = cur
	cur.right = cur

	scorev := make([]int, nplayers)
	iplayer := 0
	for num := 1; num <= nmarbles; num++ {
		l := cur.right
		r := l.right

		if num%23 != 0 {
			marble := &mring{
				left:  l,
				right: r,
				num:   num,
			}
			l.right = marble
			r.left = marble

			cur = marble
		} else {
			scorev[iplayer] += num

			for j := 0; j < 7; j++ {
				cur = cur.left
			}

			scorev[iplayer] += cur.num
			cur.left.right = cur.right
			cur.right.left = cur.left
			cur = cur.right
		}

		iplayer = (iplayer + 1) % nplayers
	}

	// find winner
	for i, s := range scorev {
		if s > score {
			winner, score = i+1, s
		}
	}

	return winner, score
}
