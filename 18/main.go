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
	var sum1, sum2 int64

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
		// Evaluate with addition and multiplication at the same precedence.
		if val, err := eval(tokens, func(tokens []int64) (int64, error) { return reduceExpr(tokens, false) }); err != nil {
			log.Fatal(err)
		} else {
			sum1 += val
		}
		// Evaluate with addition at a higher precedence than multiplication.
		if val, err := eval(tokens, func(tokens []int64) (int64, error) { return reduceExpr(tokens, true) }); err != nil {
			log.Fatal(err)
		} else {
			sum2 += val
		}
	}
	if sc.Err() != nil {
		log.Fatal(sc.Err())
	}

	fmt.Println(sum1)
	fmt.Println(sum2)
}

const (
	plus   = -1 // '+'
	times  = -2 // '*'
	lparen = -3 // '('
	rparen = -4 // ')'
)

// tokenize converts the supplied string into tokens.
// Negative values are used to represent operators and parentheses.
func tokenize(ln string) ([]int64, error) {
	var tokens []int64
	var val int64
	inVal := false

	for _, ch := range ln {
		if ch >= '0' && ch <= '9' {
			val = val*10 + int64(ch-'0')
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

// reduceFunc reduces a simple expression that doesn't contain parentheses.
type reduceFunc func(tokens []int64) (int64, error)

// reduceExpr implements reduceFunc, reducing in left-to-right order.
// If addFirst is true, addition has a higher precedence than multiplication;
// otherwise they have the same precedence.
func reduceExpr(tokens []int64, addFirst bool) (int64, error) {
	if len(tokens) == 0 || len(tokens)%2 != 1 {
		return 0, fmt.Errorf("invalid expr %v", tokens)
	}

	for pass := 0; len(tokens) > 1; pass++ {
		if pass > 1 {
			return 0, errors.New("somehow not done after two passes")
		}
		newTokens := make([]int64, 0, len(tokens))
		for i := 0; i < len(tokens); i++ {
			t := tokens[i]
			switch {
			case t >= 0:
				newTokens = append(newTokens, t)
			case t == plus:
				newTokens[len(newTokens)-1] += tokens[i+1]
				i++
			case t == times:
				if addFirst && pass == 0 {
					newTokens = append(newTokens, t) // multiply in next pass
				} else {
					newTokens[len(newTokens)-1] *= tokens[i+1]
					i++
				}
			default:
				return 0, fmt.Errorf("invalid token %v", t)
			}
		}
		tokens = newTokens
	}
	return tokens[0], nil
}

// eval evaluates tokens to a single value, using reduce to reduce
// expressions at the same depth.
func eval(tokens []int64, reduce reduceFunc) (int64, error) {
	for len(tokens) > 1 {
		// Find the depth of the most-nested expression(s).
		var depth, maxDepth int
		for _, t := range tokens {
			if t == lparen {
				if depth++; depth > maxDepth {
					maxDepth = depth
				}
			} else if t == rparen {
				depth--
			}
		}
		if depth != 0 {
			return 0, errors.New("unbalanced parentheses")
		}

		// Reduce the most-nested expression(s) to single values.
		depth = 0
		newTokens := make([]int64, 0, len(tokens))
		var exprTokens []int64
		for i := 0; i < len(tokens); i++ {
			t := tokens[i]

			switch t {
			case lparen:
				if depth++; depth == maxDepth {
					continue
				}
			case rparen:
				if depth--; depth == maxDepth-1 {
					// When we're exiting the max depth, reduce the expression.
					val, err := reduce(exprTokens)
					if err != nil {
						return 0, err
					}
					newTokens = append(newTokens, val)
					exprTokens = nil
					continue
				}
			}

			if depth == maxDepth {
				exprTokens = append(exprTokens, t) // process in this pass
			} else {
				newTokens = append(newTokens, t) // process in later pass
			}
		}

		// If we ended at the max depth (i.e. 0), reduce the final expression.
		if len(exprTokens) != 0 {
			val, err := reduce(exprTokens)
			if err != nil {
				return 0, err
			}
			newTokens = append(newTokens, val)
		}

		tokens = newTokens
	}

	return tokens[0], nil
}
