package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"strings"
)

func main() {
	black := make(map[string]struct{})
	sc := bufio.NewScanner(os.Stdin)
	for sc.Scan() {
		ln := strings.TrimSpace(sc.Text())
		if ln == "" {
			continue
		}

		var x, y int
		for i := 0; i < len(ln); i++ {
			ay := int(math.Abs(float64(y)))
			switch ln[i] {
			case 'e':
				x++
			case 'w':
				x--
			case 'n':
				switch ln[i+1] {
				case 'e':
					x += ay % 2
				case 'w':
					x -= (ay + 1) % 2
				default:
					log.Fatalf("Bad line %q", ln)
				}
				y++
				i++
			case 's':
				switch ln[i+1] {
				case 'e':
					x += ay % 2
				case 'w':
					x -= (ay + 1) % 2
				default:
					log.Fatalf("Bad line %q", ln)
				}
				y--
				i++
			default:
				log.Fatalf("Bad line %q", ln)
			}
		}

		key := fmt.Sprintf("%d|%d", x, y)
		if _, ok := black[key]; ok {
			delete(black, key) // flip to white
		} else {
			black[key] = struct{}{} // flip to black
		}
	}
	if sc.Err() != nil {
		log.Fatal(sc.Err())
	}
	fmt.Println(len(black))
}
