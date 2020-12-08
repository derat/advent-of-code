package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func main() {
	re := regexp.MustCompile(`^(\d+)-(\d+) ([a-z]): ([a-z]+)$`)
	var valid1, valid2 int
	sc := bufio.NewScanner(os.Stdin)
	for sc.Scan() {
		m := re.FindStringSubmatch(sc.Text())
		if m == nil {
			panic(fmt.Sprintf("bad line %q", sc.Text()))
		}
		min, _ := strconv.Atoi(m[1])
		max, _ := strconv.Atoi(m[2])
		ch, pw := m[3], m[4]

		if n := len(pw) - len(strings.ReplaceAll(pw, ch, "")); n >= min && n <= max {
			valid1++
		}
		if min >= 1 && min <= len(pw) && max >= 1 && max <= len(pw) &&
			((pw[min-1] == ch[0]) != (pw[max-1] == ch[0])) {
			valid2++
		}
	}
	if sc.Err() != nil {
		panic(sc.Err())
	}
	fmt.Println(valid1, valid2)
}
