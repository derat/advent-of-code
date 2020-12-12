package main

import (
	"bufio"
	"log"
	"math"
	"os"
	"strconv"
)

func main() {
	var x, y float64
	head := 90 // east

	sc := bufio.NewScanner(os.Stdin)
	for sc.Scan() {
		if sc.Text() == "" {
			continue
		}

		s := sc.Text()
		op := s[0]
		v, err := strconv.Atoi(s[1:])
		if err != nil {
			log.Fatal(err)
		}

		switch op {
		case 'N':
			y += float64(v)
		case 'S':
			y -= float64(v)
		case 'E':
			x += float64(v)
		case 'W':
			x -= float64(v)
		case 'L':
			head = (head - v) % 360
		case 'R':
			head = (head + v) % 360
		case 'F':
			// Overkill: input just uses intervals of 90 degrees for rotations.
			rad := float64(head) * math.Pi / 180
			x += math.Sin(rad) * float64(v)
			y += math.Cos(rad) * float64(v)
		default:
			log.Fatalf("invalid op %q", op)
		}
	}
	if sc.Err() != nil {
		log.Fatal(sc.Err())
	}

	println(int(math.Round(math.Abs(x) + math.Abs(y))))
}
