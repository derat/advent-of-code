package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
)

func main() {
	sum := 0
	sc := bufio.NewScanner(os.Stdin)
	for sc.Scan() {
		ln := strings.TrimSpace(sc.Text())
		if ln == "" {
			continue
		}
		tokens, err := tokenize(ln)
		if err != nil {
			log.Fatal(err)
		}
		val, err := eval(tokens)
		if err != nil {
			log.Fatal(err)
		}
		sum += val
	}
	if sc.Err() != nil {
		log.Fatal(sc.Err())
	}
	fmt.Println(sum)
}

const (
	plus   = -1
	times  = -2
	lparen = -3
	rparen = -4
)

func tokenize(ln string) ([]int, error) {
	var tokens []int
	val := 0
	inVal := false

	for _, ch := range ln {
		if ch >= '0' && ch <= '9' {
			val = val*10 + int(ch-'0')
			inVal = true
		} else {
			if inVal {
				tokens = append(tokens, val)
				val = 0
				inVal = false
			}
			switch ch {
			case '+':
				tokens = append(tokens, plus)
			case '*':
				tokens = append(tokens, times)
			case '(':
				tokens = append(tokens, lparen)
			case ')':
				tokens = append(tokens, rparen)
			case ' ':
			default:
				return nil, fmt.Errorf("bad char %q", ch)
			}
		}
	}
	if inVal {
		tokens = append(tokens, val)
	}
	return tokens, nil
}

func eval(tokens []int) (int, error) {
	// Shifts the first token off of |tokens|, reducing parenthesized expressions.
	shift := func() (int, error) {
		if len(tokens) == 0 {
			return 0, errors.New("ran out of tokens")
		}
		if tokens[0] != lparen {
			v := tokens[0]
			tokens = tokens[1:]
			return v, nil
		}

		ridx := -1
		depth := 1
		for i := 1; i < len(tokens); i++ {
			if tokens[i] == lparen {
				depth++
			} else if tokens[i] == rparen {
				depth--
				if depth == 0 {
					ridx = i
					break
				}
			}
		}
		if ridx < 0 {
			return 0, errors.New("matching rparen not found")
		}
		v, err := eval(tokens[1:ridx])
		tokens = tokens[ridx+1:]
		return v, err
	}

	for len(tokens) > 1 {
		lhs, err := shift()
		if err != nil {
			return 0, err
		}
		op, err := shift()
		if err != nil {
			return 0, err
		}
		rhs, err := shift()
		if err != nil {
			return 0, err
		}

		if lhs < 0 {
			return 0, fmt.Errorf("invalid lhs %v", lhs)
		}
		if rhs < 0 {
			return 0, fmt.Errorf("invalid rhs %v", rhs)
		}
		switch op {
		case plus:
			tokens = append([]int{lhs + rhs}, tokens...)
		case times:
			tokens = append([]int{lhs * rhs}, tokens...)
		default:
			return 0, fmt.Errorf("invalid op %v", op)
		}
	}

	if tokens[0] < 0 {
		return 0, fmt.Errorf("invalid singleton token %v", tokens[0])
	}
	return tokens[0], nil
}
