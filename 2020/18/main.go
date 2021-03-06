package main

import (
	"errors"
	"fmt"
	"log"
	"regexp"

	"github.com/derat/advent-of-code/lib"
)

func main() {
	var sum1, sum2 int64
	numRegexp := regexp.MustCompile(`^\d+`)
	for _, ln := range lib.InputLines("2020/18") {
		var tokens []int64
		for _, tok := range lib.Tokenize(ln, "+", "*", "(", ")", numRegexp) {
			switch tok {
			case "+":
				tokens = append(tokens, plus)
			case "*":
				tokens = append(tokens, times)
			case "(":
				tokens = append(tokens, lparen)
			case ")":
				tokens = append(tokens, rparen)
			default: // number
				tokens = append(tokens, lib.ExtractInt64s(tok)[0])
			}
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

	fmt.Println(sum1)
	fmt.Println(sum2)
}

const (
	plus   = -1 // '+'
	times  = -2 // '*'
	lparen = -3 // '('
	rparen = -4 // ')'
)

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
