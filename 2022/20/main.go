package main

import (
	"container/list"
	"fmt"
	"strconv"
	"strings"

	"github.com/derat/advent-of-code/lib"
)

func main() {
	input := lib.InputInts("2022/20")
	els := make([]*list.Element, len(input))
	list := list.New()
	var zero int // index of zero value
	for i, n := range input {
		els[i] = list.PushBack(n)
		if n == 0 {
			zero = i
		}
	}

	for _, el := range els {
		n := el.Value.(int)
		for n > 0 {
			if el == list.Back() {
				list.MoveAfter(el, list.Front())
			} else if el.Next() == list.Back() {
				list.MoveToFront(el)
			} else {
				list.MoveAfter(el, el.Next())
			}
			n--
		}
		for n < 0 {
			if el == list.Front() {
				list.MoveBefore(el, list.Back())
			} else if el.Prev() == list.Front() {
				list.MoveToBack(el)
			} else {
				list.MoveBefore(el, el.Prev())
			}
			n++
		}
	}

	var sum int
	el := els[zero]
	for i := 0; i <= 3000; i++ {
		if i == 1000 || i == 2000 || i == 3000 {
			sum += el.Value.(int)
		}
		el = el.Next()
		if el == nil {
			el = list.Front()
		}
	}
	fmt.Println(sum)
}

func dump(list *list.List) {
	vals := make([]string, 0, list.Len())
	for el := list.Front(); el != nil; el = el.Next() {
		vals = append(vals, strconv.Itoa(el.Value.(int)))
	}
	fmt.Println(strings.Join(vals, ", "))
}
