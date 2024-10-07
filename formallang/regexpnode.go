package formallang

import (
	"fmt"
)

type regExpNode interface {
	ToString(priority int) string
	Priority() int
}

type regExpNodeRune struct {
	r rune
}

const (
	lowPriority = iota
	addPriority
	mulPriority
	cliniPriority
	runePriority
	hightPriority
)

func (regExpNodeRune) Priority() int         { return runePriority }
func (r regExpNodeRune) ToString(int) string { return fmt.Sprintf("%c", r.r) }

type regExpNodeAdd struct {
	Next0, Next1 regExpNode
}

func (regExpNodeAdd) Priority() int { return addPriority }
func (add regExpNodeAdd) ToString(priority int) string {
	prior := add.Priority()

	format := "%v + %v"
	if prior < priority {
		format = "(%v + %v)"
	}

	return fmt.Sprintf(format, add.Next0.ToString(prior), add.Next1.ToString(prior))
}

type regExpNodeMul struct {
	Next0, Next1 regExpNode
}

func (regExpNodeMul) Priority() int { return mulPriority }
func (mul regExpNodeMul) ToString(priority int) string {
	prior := mul.Priority()

	format := "%v%v"
	if prior < priority {
		format = "(%v%v)"
	}

	return fmt.Sprintf(format, mul.Next0.ToString(prior), mul.Next1.ToString(prior))
}

type regExpNodeClini struct {
	Next regExpNode
}
func (regExpNodeClini) Priority() int { return cliniPriority }
func (clini regExpNodeClini) ToString(priority int) string {
	prior := clini.Priority()
	return fmt.Sprintf("%v*", clini.Next.ToString(prior))
}

func createRegExpNodes(tokens []Token) (regExpNode, error) {
	start := 0

	res, err := createRegExpNodesRecursive(tokens, &start, len(tokens))

	if err == nil && start < len(tokens) {
		return nil, fmt.Errorf("error at %v idx symbol", start)
	}

	return res, err
}

func createRegExpNodesRecursiveBrases(tokens []Token, idx *int) (regExpNode, error) {
	if *idx >= len(tokens) {
		return nil, fmt.Errorf("can't parse")
	}
	start := *idx

	for {
		if !(tokens[*idx].Symb == '(' && tokens[*idx].Servicable) {
			break
		}
		(*idx)++

		res, err := createRegExpNodesRecursiveSum(tokens, idx, potentialEnd)
		if err != nil {
			break
		}

		if *idx >= len(tokens) {
			break
		}
		if !(tokens[*idx].Symb == ')' && tokens[*idx].Servicable) {
			break
		}
		(*idx)++

		return res, nil
	}

	*idx = start
	return nil, fmt.Errorf("error at %v idx symbol", *idx)
}

func createRegExpNodesRecursiveSum(tokens []Token, idx *int) (regExpNode, error) {
	if *idx >= len(tokens) {
		return nil, fmt.Errorf("can't parse")
	}
	start := *idx

loop:
	for {
		res, err := createRegExpNodesRecursiveMul(tokens, idx)
		if err != nil {
			break
		}

		ptr := &res
		for (*idx+1 < len(tokens) && tokens[*idx] == Token{Symb: '+', Servicable: true}) {
			buf, err := createRegExpNodesRecursiveMul(tokens, idx)
			if err != nil {
				break loop
			}

			newAddNode := regExpNodeAdd{*ptr, buf}
			ptr = &newAddNode.Next1
		}

		return res, nil
	}

	*idx = start
	return nil, fmt.Errorf("error at %v idx symbol", *idx)
}

func createRegExpNodesRecursiveMul(tokens []Token, idx *int) (regExpNode, error) {
	if *idx >= len(tokens) {
		return nil, fmt.Errorf("can't parse")
	}
	start := *idx

loop:
	for {
		for {
			res, err := createRegExpNodesRecursiveMul(tokens, idx)
			if err != nil {
				break
			}

			res, er
		}
		

		ptr := &res
		for 

		return res, nil
	}

	*idx = start
	return nil, fmt.Errorf("error at %v idx symbol", *idx)
}

func createRegExpNodesRecursive(tokens []Token, idx *int, potentialEnd int) (regExpNode, error) {
	if *idx >= potentialEnd {
		return nil, fmt.Errorf("can't parse")
	}

	start := *idx


	/*** expr = expr* ***/
	*idx = start
	for {
		Next, err := createRegExpNodesRecursive(tokens, idx, potentialEnd-1)
		if err != nil {
			break
		}

		if *idx >= len(tokens) {
			break
		}
		if !(tokens[*idx].Symb == '*' && tokens[*idx].Servicable) {
			break
		}
		(*idx)++

		return regExpNodeClini{Next}, nil
	}

	*idx = start
	/*** expr = expr+ ***/
	for {
		Next, err := createRegExpNodesRecursive(tokens, idx, potentialEnd-1)
		if err != nil {
			break
		}

		if *idx >= potentialEnd {
			break
		}
		if !(tokens[*idx].Symb == '*' && tokens[*idx].Servicable) {
			break
		}
		(*idx)++

		// TODO copy
		return regExpNodeAdd{Next, regExpNodeClini{Next}}, nil
	}

	*idx = start
	/*** expr = rune ***/
	for {
		if tokens[*idx].Servicable {
			break
		}
		bufidx := *idx
		(*idx)++
		return regExpNodeRune{tokens[bufidx].Symb}, nil
	}

	return nil, fmt.Errorf("error at %v idx symbol", *idx)
}
