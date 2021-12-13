package main

import (
	"fmt"
	"log"
	"os"
	"runtime/pprof"
	"sort"
	"strings"
	"unicode"

	"github.com/derat/advent-of-code/lib"
)

const (
	prof    = false
	version = 3
)

func main() {
	if prof {
		f, err := os.Create("cpu.prof")
		if err != nil {
			log.Fatal("Could not create CPU profile: ", err)
		}
		defer f.Close()
		if err := pprof.StartCPUProfile(f); err != nil {
			log.Fatal("Could not start CPU profile: ", err)
		}
		defer pprof.StopCPUProfile()
	}

	conns := make(map[string][]string)
	for _, ln := range lib.InputLines("2021/12") {
		parts := strings.Split(ln, "-")
		lib.AssertEq(len(parts), 2)
		src, dst := parts[0], parts[1]
		// A connection between two big caves would result in an infinite number of paths.
		lib.Assert(!(unicode.IsUpper(rune(src[0])) && unicode.IsUpper(rune(dst[0]))))
		conns[src] = append(conns[src], dst)
		conns[dst] = append(conns[dst], src)
	}

	var f func(map[string][]string, bool) int
	switch version {
	case 1:
		f = exploreSlow
	case 2:
		f = exploreFast
	case 3:
		f = exploreFaster
	default:
		panic("Invalid version")
	}

	// Part 1: "Your goal is to find the number of distinct paths that start at start, end at end,
	// and don't visit small caves more than once. There are two types of caves: big caves (written
	// in uppercase, like A) and small caves (written in lowercase, like b). It would be a waste of
	// time to visit any small cave more than once, but big caves are large enough that it might be
	// worth visiting them multiple times. So, all paths you find should visit small caves at most
	// once, and can visit big caves any number of times."
	fmt.Println(f(conns, false))

	// Part 2: "After reviewing the available paths, you realize you might have time to visit a
	// single small cave twice. Specifically, big caves can be visited any number of times, a single
	// small cave can be visited at most twice, and the remaining small caves can be visited at most
	// once. However, the caves named start and end can only be visited exactly once each: once you
	// leave the start cave, you may not return to it, and once you reach the end cave, the path
	// must end immediately."
	fmt.Println(f(conns, true))
}

// This is the first version I wrote. It took about 3s, and profiling showed it spending
// a bunch of time hashing strings. After adding NoSteps/NoFrom, it went down to about 1.5s.
func exploreSlow(conns map[string][]string, canVisitSmallTwice bool) (npaths int) {
	lib.BFS([]interface{}{"start"}, func(si interface{}, next map[interface{}]struct{}) {
		s := si.(string)
		if strings.HasSuffix(s, ",end") {
			npaths++
			return
		}

		path := strings.Split(s, ",")                // caves in visited order
		seen := make(map[string]struct{}, len(path)) // set of visited caves
		twice := false
		for _, cave := range path {
			if _, ok := seen[cave]; ok && unicode.IsLower(rune(cave[0])) {
				twice = true
			}
			seen[cave] = struct{}{}
		}

		at := path[len(path)-1]
		for _, dst := range conns[at] {
			if _, ok := seen[dst]; !ok || unicode.IsUpper(rune(dst[0])) ||
				(canVisitSmallTwice && !twice && dst != "start" && dst != "end") {
				next[s+","+dst] = struct{}{}
			}
		}
	}, &lib.BFSOptions{NoSteps: true, NoFrom: true})

	return npaths
}

// This is the second version. I wasted way too much time trying to pack the full state
// into a uint64 before realizing that I still need to track the order in which caves are
// visited. It took about 2.3s, but then went down to about 0.6s with NoSteps/NoFrom.
func exploreFast(conns map[string][]string, canVisitSmallTwice bool) (npaths int) {
	caves := lib.MapStringKeys(conns)
	sort.Strings(caves)
	caveNum := make(map[string]int, len(caves))
	for i, cave := range caves {
		caveNum[cave] = i
	}
	startNum, endNum := caveNum["start"], caveNum["end"]

	start := state{startNum, 1 << startNum, false}
	lib.BFS([]interface{}{start}, func(si interface{}, next map[interface{}]struct{}) {
		st := si.(state)
		if lib.HasBit(st.seen, endNum) {
			npaths++
			return
		}
		for _, dst := range conns[caves[st.at]] {
			big := unicode.IsUpper(rune(dst[0]))
			num := caveNum[dst]
			seen := lib.HasBit(st.seen, num)
			if big || !seen || (canVisitSmallTwice && !st.twice && dst != "start" && dst != "end") {
				ns := state{
					at:    num,
					seen:  lib.SetBit(st.seen, num, true),
					twice: st.twice || (!big && seen),
				}
				next[ns] = struct{}{}
			}
		}
	}, &lib.BFSOptions{NoSteps: true, NoFrom: true})
	return npaths
}

type state struct {
	at    int    // cave number
	seen  uint64 // bitfield of visited caves from caveNum
	twice bool   // a small cave has already been visited twice
}

// This is the third version. I realized that since there are no connections between big in the
// input (since that would result in an infinite number of paths), I can actually use a bitfield
// after all as long the BFS function is told to not check for duplicate states itself. This brings
// the running time down to 0.5s, probably due to avoiding exploreFast's struct allocs.
func exploreFaster(conns map[string][]string, canVisitSmallTwice bool) (npaths int) {
	caves := lib.MapStringKeys(conns)
	sort.Strings(caves)
	caveNum := make(map[string]int, len(caves))
	for i, cave := range caves {
		caveNum[cave] = i
	}
	startNum, endNum := caveNum["start"], caveNum["end"]

	start := pack(startNum, 1<<startNum, false)
	lib.BFS([]interface{}{start}, func(si interface{}, next map[interface{}]struct{}) {
		at, seen, twice := unpack(si.(uint64))
		if lib.HasBit(seen, endNum) {
			npaths++
			return
		}
		for _, dst := range conns[caves[at]] {
			big := unicode.IsUpper(rune(dst[0]))
			num := caveNum[dst]
			vis := lib.HasBit(seen, num)
			if big || !vis || (canVisitSmallTwice && !twice && dst != "start" && dst != "end") {
				ns := pack(num, lib.SetBit(seen, num, true), twice || (!big && vis))
				next[ns] = struct{}{}
			}
		}
	}, &lib.BFSOptions{NoSteps: true, NoFrom: true})
	return npaths
}

func pack(at int, seen uint64, twice bool) uint64 {
	lib.Assert(seen < 1<<32)
	s := lib.PackInts(at, int(seen))
	return lib.SetBit(s, 63, twice)
}

func unpack(s uint64) (at int, seen uint64, twice bool) {
	twice = lib.HasBit(s, 63)
	lib.SetBit(s, 63, false)
	at, si := lib.UnpackInt2(s)
	return at, uint64(si), twice
}
