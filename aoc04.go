package main

import (
	"fmt"
	"sort"
	"time"

	"github.com/pkg/errors"
)

type WatchEntryType int

const (
	BeginShift WatchEntryType = iota
	FallsAsleep
	WakesUp
)

type WatchEntry struct {
	Time time.Time
	Type WatchEntryType
	ID   int // only when BeginShift
}

func ParseOnWatchLog(log []string) ([]WatchEntry, error) {
	var v []WatchEntry
	for _, s := range log {
		e, err := ParseOnWatchEntry(s)
		if err != nil {
			return nil, err
		}
		v = append(v, e)
	}
	sort.Slice(v, func(i, j int) bool {
		return v[i].Time.Before(v[j].Time)
	})
	return v, nil
}

const watchEntryTimeFmt = "2006-01-02 15:04"

func (we *WatchEntry) String() string {
	ts := we.Time.Format(watchEntryTimeFmt)
	var msg string
	switch we.Type {
	case BeginShift:
		msg = fmt.Sprintf("Guard #%d begins shift", we.ID)
	case FallsAsleep:
		msg = "falls asleep"
	case WakesUp:
		msg = "wakes up"
	default:
		msg = "invalid"
	}

	return fmt.Sprintf("[%s] %s", ts, msg)
}

/* Format:
[1518-04-10 00:00] Guard #2819 begins shift
[1518-10-01 00:56] wakes up
[1518-09-10 00:52] wakes up
[1518-08-14 00:53] falls asleep
*/

func ParseOnWatchEntry(s string) (WatchEntry, error) {
	if len(s) < 21 || s[0] != '[' || s[17:19] != "] " {
		return WatchEntry{}, errors.Errorf("invalid format: %v", s)
	}
	ts := s[1:17]
	msg := s[19:]
	t, err := time.Parse(watchEntryTimeFmt, ts)
	if err != nil {
		return WatchEntry{}, errors.Errorf("invalid timestamp in %v", s)
	}

	e := WatchEntry{
		Time: t,
	}
	switch {
	case msg == "falls asleep":
		e.Type = FallsAsleep
	case msg == "wakes up":
		e.Type = WakesUp
	default:
		_, err := fmt.Sscanf(msg, "Guard #%d begins shift", &e.ID)
		if err != nil {
			return WatchEntry{}, errors.Errorf("unrecognised message in %v", s)
		}
		e.Type = BeginShift
	}
	return e, nil
}

func FindMostSleepyGuard(log []WatchEntry) (id, minute int) {
	gm := mapWatchLog(log)

	xid, xm, xasleep := 0, 0, 0
	for id, e := range gm {
		if e.asleep > xasleep {
			xid = id
			xasleep = e.asleep

			// find best minute
			maxm, maxcount := 0, 0
			for m, count := range e.chart {
				if count > maxcount {
					maxm, maxcount = m, count
				}
			}

			xm = maxm
		}
	}
	return xid, xm
}

func FindMostSleptMinute(log []WatchEntry) (id, minute int) {
	gm := mapWatchLog(log)

	xid, xm, xcount := 0, 0, 0
	for id, e := range gm {
		// find best minute
		maxm, maxcount := 0, 0
		for m, count := range e.chart {
			if count > maxcount {
				maxm, maxcount = m, count
			}
		}

		if maxcount > xcount {
			xid = id
			xm = maxm
			xcount = maxcount
		}
	}
	return xid, xm
}

type guardWatchEntry struct {
	chart  [60]int // minute chart
	asleep int     // total minutes slept
}

func mapWatchLog(log []WatchEntry) map[int]guardWatchEntry {
	gm := make(map[int]guardWatchEntry)
	id := 0
	asleep := false
	masleep := 0
	for _, we := range log {
		switch we.Type {
		case BeginShift:
			if asleep {
				panic("guard still asleep")
			}
			id, asleep = we.ID, false
		case FallsAsleep:
			if asleep {
				panic("sleeping guard falls asleep")
			}
			masleep, asleep = we.Time.Minute(), true
		case WakesUp:
			if !asleep {
				panic("awake guard wakes up")
			}
			mwakeup := we.Time.Minute()
			e := gm[id]
			e.asleep += (mwakeup - masleep)
			for m := masleep; m < mwakeup; m++ {
				e.chart[m]++
			}
			gm[id] = e
			asleep = false
		}
	}
	return gm
}
