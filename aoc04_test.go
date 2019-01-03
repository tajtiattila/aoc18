package main

import (
	"strings"
	"testing"
)

func TestAoC04(t *testing.T) {
	data := input04(t)
	id, m := FindMostSleepyGuard(data)
	r := id * m
	t.Logf("#1 id=%d minute=%d; %d×%d=%d", id, m, id, m, id*m)
	// 19896 125362
	if r <= 125362 {
		t.Fatal("too low")
	}

	id, m = FindMostSleptMinute(data)
	r = id * m
	t.Logf("#2 id=%d minute=%d; %d×%d=%d", id, m, id, m, id*m)
}

func input04(t *testing.T) []WatchEntry {
	r, err := ParseOnWatchLog(input04v)
	if err != nil {
		t.Fatalf("parse on watch log: %v", err)
	}
	return r
}

func TestAoC04Sample(t *testing.T) {
	data, err := ParseOnWatchLog(sample04v)
	if err != nil {
		t.Fatalf("parse on watch log: %v", err)
	}
	id, m := FindMostSleepyGuard(data)
	if id != 10 || m != 24 {
		t.Fatalf("FindMostSleepyGuard got id=%d, minute=%d; want id=10, minute=24", id, m)
	}
	id, m = FindMostSleptMinute(data)
	if id != 99 || m != 45 {
		t.Fatalf("FindMostSleptMinute got id=%d, minute=%d; want id=99, minute=45", id, m)
	}
}

var sample04v = strings.Split(`[1518-11-01 00:00] Guard #10 begins shift
[1518-11-01 00:05] falls asleep
[1518-11-01 00:25] wakes up
[1518-11-01 00:30] falls asleep
[1518-11-01 00:55] wakes up
[1518-11-01 23:58] Guard #99 begins shift
[1518-11-02 00:40] falls asleep
[1518-11-02 00:50] wakes up
[1518-11-03 00:05] Guard #10 begins shift
[1518-11-03 00:24] falls asleep
[1518-11-03 00:29] wakes up
[1518-11-04 00:02] Guard #99 begins shift
[1518-11-04 00:36] falls asleep
[1518-11-04 00:46] wakes up
[1518-11-05 00:03] Guard #99 begins shift
[1518-11-05 00:45] falls asleep
[1518-11-05 00:55] wakes up`, "\n")
