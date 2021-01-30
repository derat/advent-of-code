package main

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/derat/advent-of-code/lib"
)

func main() {
	var nvalid, nvalid2 int
	for _, pg := range lib.InputParagraphs("2020/4") {
		data := strings.Join(pg, " ")
		if valid(data, false) {
			nvalid++
		}
		if valid(data, true) {
			nvalid2++
		}
	}
	fmt.Println(nvalid)
	fmt.Println(nvalid2)
}

var req = map[string]*regexp.Regexp{
	"byr": regexp.MustCompile(`^(19[2-9]\d|200[0-2])$`),                     // (Birth Year)
	"iyr": regexp.MustCompile(`^(201\d|2020)$`),                             // (Issue Year)
	"eyr": regexp.MustCompile(`^(202\d|2030)$`),                             // (Expiration Year)
	"hgt": regexp.MustCompile(`^((1[5-8]\d|19[0-3])cm|(59|6\d|7[0-6])in)$`), // (Height)
	"hcl": regexp.MustCompile(`^#[0-9a-f]{6}$`),                             // (Hair Color)
	"ecl": regexp.MustCompile(`^(amb|blu|brn|gry|grn|hzl|oth)$`),            // (Eye Color)
	"pid": regexp.MustCompile(`^\d{9}$`),                                    // (Passport ID)
}

var opt = map[string]struct{}{
	"cid": struct{}{}, // (Country ID)
}

func valid(s string, check bool) bool {
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

	for key, re := range req {
		val := fields[key]
		if val == "" {
			return false
		}
		if check && !re.MatchString(val) {
			return false
		}
	}
	return true
}
