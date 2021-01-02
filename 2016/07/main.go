package main

import (
	"fmt"
	"regexp"

	"github.com/derat/advent-of-code/lib"
)

func main() {
	bareRegexp := regexp.MustCompile(`^[a-z]+`)
	netRegexp := regexp.MustCompile(`^\[[a-z]+\]`)
	var tlsCnt, sslCnt int
	for _, ln := range lib.InputLines("2016/7") {
		var bareCnt, netCnt int
		bareSet := make(map[string]struct{})
		netSet := make(map[string]struct{})
		for _, tok := range lib.Tokenize(ln, bareRegexp, netRegexp) {
			var net bool
			if tok[0] == '[' {
				net = true
				tok = tok[1 : len(tok)-1]
			}
			for i := 0; i < len(tok)-3; i++ {
				if tok[i] != tok[i+1] && tok[i] == tok[i+3] && tok[i+1] == tok[i+2] {
					if net {
						netCnt++
					} else {
						bareCnt++
					}
					break
				}
			}
			for i := 0; i < len(tok)-2; i++ {
				if tok[i] != tok[i+1] && tok[i] == tok[i+2] {
					if net {
						netSet[tok[i:i+3]] = struct{}{}
					} else {
						bareSet[tok[i:i+3]] = struct{}{}
					}
				}
			}
		}
		if bareCnt > 0 && netCnt == 0 {
			tlsCnt++
		}
		for s := range bareSet {
			rev := s[1:2] + s[0:1] + s[1:2]
			if _, ok := netSet[rev]; ok {
				sslCnt++
				break
			}
		}
	}
	fmt.Println(tlsCnt)
	fmt.Println(sslCnt)
}
