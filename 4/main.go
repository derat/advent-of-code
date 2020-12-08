package main

import (
	"bufio"
	"os"
	"strings"
)

func main() {
	sc := bufio.NewScanner(os.Stdin)
	var nvalid int
	var data string
	for sc.Scan() {
		if sc.Text() != "" {
			data += " " + sc.Text()
		} else {
			if valid(data) {
				nvalid++
			}
			data = ""
			continue
		}
	}
	if sc.Err() != nil {
		panic(sc.Err())
	}
	if valid(data) {
		nvalid++
	}
	println(nvalid)
}

var req = map[string]struct{}{
	"byr": struct{}{}, // (Birth Year)
	"iyr": struct{}{}, // (Issue Year)
	"eyr": struct{}{}, // (Expiration Year)
	"hgt": struct{}{}, // (Height)
	"hcl": struct{}{}, // (Hair Color)
	"ecl": struct{}{}, // (Eye Color)
	"pid": struct{}{}, // (Passport ID)
}

var opt = map[string]struct{}{
	"cid": struct{}{}, // (Country ID)
}

func valid(s string) bool {
	fields := make(map[string]string)
	for _, x := range strings.Fields(s) {
		p := strings.SplitN(x, ":", 2)
		if len(p) != 2 {
			return false
		}
		key, val := p[0], p[1]
		fields[key] = val
		// TODO: Instructions don't seem to say anything about unexpected fields.
	}

	for key := range req {
		if fields[key] == "" {
			return false
		}
	}
	return true
}
