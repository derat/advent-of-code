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
	valid := 0
	sc := bufio.NewScanner(os.Stdin)
	for sc.Scan() {
		m := re.FindStringSubmatch(sc.Text())
		if m == nil {
			panic(fmt.Sprintf("bad line %q", sc.Text()))
		}
		min, _ := strconv.Atoi(m[1])
		max, _ := strconv.Atoi(m[2])
		ch, pw := m[3], m[4]
		n := len(pw) - len(strings.ReplaceAll(pw, ch, ""))
		if n >= min && n <= max {
			valid++
		}
	}
	if sc.Err() != nil {
		panic(sc.Err())
	}
	fmt.Println(valid)
}
