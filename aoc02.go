package main

import "strings"

func BoxIDListChecksum(ids []string) int {
	c2, c3 := 0, 0
	for _, id := range ids {
		i2, i3 := boxSum(id, 2), boxSum(id, 3)
		c2 += i2
		c3 += i3
	}
	return c2 * c3
}

func boxSum(id string, n int) int {
	count := make(map[rune]int)
	for _, c := range id {
		count[c]++
	}
	for _, v := range count {
		if v == n {
			return 1
		}
	}
	return 0
}

func BoxIDListSimilar(ids []string) []string {
	matcher := make(map[string][]string)
	for _, id := range ids {
		for _, m := range boxIdSimMatchers(id) {
			matcher[m] = append(matcher[m], id)
		}
	}

	var res []string
	for m, v := range matcher {
		if len(v) > 1 {
			res = append(res, strings.Replace(m, "?", "", 1))
		}
	}
	return res
}

func boxIdSimMatchers(id string) []string {
	var v []string
	p := []byte(id)
	for i, c := range p {
		p[i] = '?'
		v = append(v, string(p))
		p[i] = c
	}
	return v
}
