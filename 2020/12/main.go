package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
)

func main() {
	// Part 1:
	var ox, oy float64    // ship coords
	var head float64 = 90 // ship heading (east)

	// Part 2:
	var sx, sy float64         // ship coords
	var wx, wy float64 = 10, 1 // waypoint relative to ship

	rotWay := func(deg float64) {
		dist := math.Sqrt(math.Pow(wx, 2) + math.Pow(wy, 2))
		// atan2 is needed to preserve the signs of the individual coordinates:
		// when passing a ratio to atan, (10,1) and (-10,-1) are indistinguishable.
		rad := math.Atan2(wx, wy) + (deg * math.Pi / 180)
		wx = math.Sin(rad) * dist
		wy = math.Cos(rad) * dist
	}

	sc := bufio.NewScanner(os.Stdin)
	for sc.Scan() {
		if sc.Text() == "" {
			continue
		}

		s := sc.Text()
		op := s[0]
		v, err := strconv.ParseFloat(s[1:], 64)
		if err != nil {
			log.Fatal(err)
		}

		switch op {
		case 'N':
			oy += v
			wy += v
		case 'S':
			oy -= v
			wy -= v
		case 'E':
			ox += v
			wx += v
		case 'W':
			ox -= v
			wx -= v
		case 'L':
			head = math.Mod(head-v, 360)
			rotWay(-v)
		case 'R':
			head = math.Mod(head+v, 360)
			rotWay(v)
		case 'F':
			// Overkill: input just uses intervals of 90 degrees for rotations.
			rad := head * math.Pi / 180
			ox += math.Sin(rad) * v
			oy += math.Cos(rad) * v

			sx += wx * v
			sy += wy * v
		default:
			log.Fatalf("invalid op %q", op)
		}
	}
	if sc.Err() != nil {
		log.Fatal(sc.Err())
	}

	// https://en.wikipedia.org/wiki/Taxicab_geometry
	fmt.Println(int(math.Round(math.Abs(ox) + math.Abs(oy))))
	fmt.Println(int(math.Round(math.Abs(sx) + math.Abs(sy))))
}
