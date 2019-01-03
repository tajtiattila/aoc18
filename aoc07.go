package main

import (
	"fmt"
	"sort"
)

type AssemblyLink struct {
	Pred, Succ string
}

func ParseAssemblyLink(s string) (a AssemblyLink, err error) {
	_, err = fmt.Sscanf(s, "Step %s must be finished before step %s can begin.", &a.Pred, &a.Succ)
	return
}

func SortAssemblyInstr(links []AssemblyLink) []string {
	inf := prepAssembly(links)

	var steps []string

	for !inf.allDone() {
		inst, ok := inf.pickWork()
		if !ok {
			panic("logic error")
		}

		steps = append(steps, inst)
		inf.markDone(inst)
	}

	return steps
}

func TimeAssembly(links []AssemblyLink, numworkers int, instrtimef func(string) int) int {
	if numworkers < 1 {
		panic("numworkers must be positive")
	}
	inf := prepAssembly(links)

	type worker struct {
		inst string // current instruction
		work int    // time to finish
	}

	var busy []worker
	nfree := numworkers

	t := 0

	for !inf.allDone() {

		// pick work for available workers
		for nfree > 0 {
			inst, ok := inf.pickWork()
			if !ok {
				break // no work available
			}

			nfree--
			busy = append(busy, worker{
				inst: inst,
				work: instrtimef(inst),
			})
		}

		// process work
		var nbusy int
		for i := range busy {
			w := &busy[i]
			w.work--
			if w.work > 0 {
				busy[nbusy] = *w
				nbusy++
			} else {
				nfree++
				inf.markDone(w.inst)
			}
		}
		busy = busy[:nbusy]

		t++
	}

	return t
}

const (
	statusNone   = 0
	statusPicked = 1
	statusDone   = 2
)

type assemblyInfo struct {
	names []string

	m map[string]*assemblyEnt

	npickable int // number of instructions available (not picked)
	ndoable   int // number of instructions not done yet
}

type assemblyEnt struct {
	pred []string

	status int
}

func prepAssembly(links []AssemblyLink) *assemblyInfo {
	allnames := make(map[string]struct{})
	for _, link := range links {
		allnames[link.Pred] = struct{}{}
		allnames[link.Succ] = struct{}{}
	}

	m := make(map[string]*assemblyEnt, len(allnames))
	for n := range allnames {
		m[n] = &assemblyEnt{}
	}

	for _, link := range links {
		x := m[link.Succ]
		x.pred = append(x.pred, link.Pred)
	}

	var names []string
	for n := range allnames {
		names = append(names, n)
	}
	sort.Strings(names)

	return &assemblyInfo{
		names: names,
		m:     m,

		npickable: len(names),
		ndoable:   len(names),
	}
}

func (inf *assemblyInfo) allDone() bool {
	return inf.ndoable == 0
}

func (inf *assemblyInfo) pickWork() (inst string, ok bool) {
	for _, inst := range inf.names {
		x := inf.m[inst]
		if x.status != statusNone {
			continue // already started
		}

		able := true
		for _, pred := range x.pred {
			if inf.m[pred].status != statusDone {
				able = false
				break // pred still unfinished
			}
		}

		if able {
			x.status = statusPicked
			inf.npickable--
			return inst, true
		}
	}

	return "", false
}

func (inf *assemblyInfo) markDone(inst string) {
	x := inf.m[inst]
	if x == nil || x.status != statusPicked {
		panic("logic error")
	}
	x.status = statusDone
	inf.ndoable--
}
