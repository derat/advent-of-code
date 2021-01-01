package main

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"strconv"

	"github.com/derat/advent-of-code/lib"
)

func main() {
	const (
		zeros = 5 // number of leading zeros
		pwLen = 8
	)
	pre := lib.InputLines("2016/5")[0]
	want := make([]byte, zeros/2) // full zero bytes
	var pw string
	for i := 0; len(pw) < pwLen; i++ {
		s := pre + strconv.Itoa(i)
		hash := md5.Sum([]byte(s))
		if bytes.HasPrefix(hash[:], want) && (zeros%2 == 0 || hash[zeros/2] < 16) {
			b := lib.IfByte(zeros%2 == 0, hash[zeros/2]>>4, hash[zeros/2]) & 0xf
			ch := fmt.Sprintf("%x", b)
			pw += ch
			fmt.Print(ch) // print a byte at a time because it looks cooler
		}
	}
	fmt.Println()
}
