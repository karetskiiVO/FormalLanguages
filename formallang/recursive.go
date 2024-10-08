package formallang

import (
	"fmt"
	//"reflect"
)

func createRegExpNodes(tokens []Token) (regExpNode, error) {
	var start = 0

	res, err := recursiveGetSum(tokens, &start)
	if err != nil {
		return nil, err
	}

	if start != len(tokens) {
		return nil, fmt.Errorf("can't parse on index %v", start)
	}

	return res, err
}

func recursiveGetRune(tokens []Token, idx *int) (regExpNode, error) {
	if *idx >= len(tokens) {
		return nil, fmt.Errorf("can't parse on index %v", *idx)
	}

	var res regExpNode
	content := tokens[*idx]
	if content.Servicable {
		if content.Symb != '1' {
			return nil, fmt.Errorf("can't parse on index %v", *idx)
		}

		res = regExpNodeEmptyRune{}
	} else {
		res = regExpNodeRune{content.Symb}
	}

	(*idx)++
	return res, nil
}

func recursiveGetBrasClini(tokens []Token, idx *int) (regExpNode, error) {
	if *idx >= len(tokens) {
		return nil, fmt.Errorf("can't parse on index %v", *idx)
	}
	start := *idx

	for {
		var res regExpNode
		var err error
		if tokens[*idx].Servicable && tokens[*idx].Symb == '(' {
			(*idx)++

			res, err = recursiveGetSum(tokens, idx)
			if err != nil {
				break
			}

			if *idx >= len(tokens) {
				break
			}
			if !(tokens[*idx].Servicable && tokens[*idx].Symb == ')') {
				break
			}
			(*idx)++
		} else {
			res, err = recursiveGetRune(tokens, idx)
			if err != nil {
				break
			}
		}

		for (*idx) < len(tokens) && tokens[*idx].Servicable && tokens[*idx].Symb == '*' {
			res = regExpNodeClini{res}
			(*idx)++
		}

		return res, nil
	}

	*idx = start
	return nil, fmt.Errorf("can't parse on index %v", *idx)
}

func recursiveGetSum(tokens []Token, idx *int) (regExpNode, error) {
	if *idx >= len(tokens) {
		return nil, fmt.Errorf("can't parse on index %v", *idx)
	}
	start := *idx

loop:
	for {
		res, err := recursiveGetMul(tokens, idx)
		if err != nil {
			break
		}

		nodes := make([]regExpNode, 1)
		nodes[0] = res

		for *idx < len(tokens) && tokens[*idx].Symb == '+' {
			(*idx)++
			buf, err := recursiveGetMul(tokens, idx)
			if err != nil {
				break loop
			}

			nodes = append(nodes, buf)
		}

		if len(nodes) > 1 {
			return regExpNodeAdd{nodes}, nil	
		}

		return nodes[0], nil
	}

	*idx = start
	return nil, fmt.Errorf("can't parse on index %v", *idx)
}

func recursiveGetMul(tokens []Token, idx *int) (regExpNode, error) {
	if *idx >= len(tokens) {
		return nil, fmt.Errorf("can't parse on index %v", *idx)
	}
	start := *idx

	for {
		res, err := recursiveGetBrasClini(tokens, idx)
		if err != nil {
			break
		}
		
		nodes := make([]regExpNode, 1)
		nodes[0] = res

		for {
			buf, err := recursiveGetBrasClini(tokens, idx)
			if err != nil {
				break
			}
			
			nodes = append(nodes, buf)
		}

		if len(nodes) > 1 {
			return regExpNodeMul{nodes}, nil	
		}

		return nodes[0], nil
	}

	*idx = start
	return nil, fmt.Errorf("can't parse on index %v", *idx)
}
