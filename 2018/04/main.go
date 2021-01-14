package main

import (
	"fmt"
	"math/bits"
	"sort"

	"github.com/derat/advent-of-code/lib"
)

func main() {
	lines := lib.InputLines("2018/4")
	sort.Strings(lines)

	var guard int
	var day uint64
	data := make(map[int][]uint64) // days keyed by guard ID
	for _, ln := range lines {
		var date, rest string
		var h, m int
		lib.Extract(ln, `^\[(\d{4}-\d\d-\d\d) (\d\d):(\d\d)\] (.+)`, &date, &h, &m, &rest)

		var newGuard int
		switch {
		case lib.ExtractMaybe(rest, `^Guard #(\d+) begins shift$`, &newGuard):
			if guard > 0 {
				data[guard] = append(data[guard], day)
			}
			guard = newGuard
			day = 0
		case rest == "falls asleep" || rest == "wakes up":
			val := rest == "falls asleep"
			for i := m; i < 60; i++ {
				day = lib.SetBit(day, i, val)
			}
		default:
			lib.Panicf("Bad line %q", ln)
		}
	}
	if guard > 0 {
		data[guard] = append(data[guard], day)
	}

	// Part 1: Print ID of most-asleep guard multiplied by their most-asleep minute.
	var worst, maxSleep int
	for guard, days := range data {
		var sum int
		for _, day := range days {
			sum += bits.OnesCount64(day)
		}
		if sum > maxSleep {
			worst = guard
			maxSleep = sum
		}
	}
	min, _ := worstMinute(data[worst])
	fmt.Println(worst * min)

	// Part 2: Print ID of guard most asleep in a single minute times that minute.
	var worst2, min2, maxCnt int
	for guard, days := range data {
		if min, cnt := worstMinute(days); cnt > maxCnt {
			worst2 = guard
			min2 = min
			maxCnt = cnt
		}
	}
	fmt.Println(worst2 * min2)
}

func worstMinute(days []uint64) (min, cnt int) {
	var mins [60]int
	for _, day := range days {
		for i := 0; i < 60; i++ {
			if lib.HasBit(day, i) {
				mins[i]++
			}
		}
	}
	var max int
	for i, cnt := range mins {
		if cnt > max {
			min = i
			max = cnt
		}
	}
	return min, max
}
