package formallang

type stack []any

// Solve7 solution for 7yh task
func Solve7(pol string, x rune, k int) bool {
	var nodes stack

	for _, r := range pol {
		switch r {
		case x:
			nodes.Push(&regExpNodeRune{r})
		case '+':
			rv := nodes.Top().(regExpNode)
			nodes.Pop()
			lv := nodes.Top().(regExpNode)
			nodes.Pop()

			nodes.Push(&regExpNodeAdd{
				[]regExpNode{lv, rv},
			})
		case '*':
			val := nodes.Top().(regExpNode)
			nodes.Pop()

			nodes.Push(&regExpNodeClini{val})
		case '.':
			rv := nodes.Top().(regExpNode)
			nodes.Pop()
			lv := nodes.Top().(regExpNode)
			nodes.Pop()

			nodes.Push(&regExpNodeMul{
				[]regExpNode{lv, rv},
			})
		default:
			nodes.Push(&regExpNodeEmptyRune{})
		}
	}

	reg := &RegExp{
		abc:  map[rune]struct{}{x: struct{}{}},
		tree: nodes.Top().(regExpNode),
	}

	dfa := DFAfromNFA(NFAFromRegExp(reg).RemoveEmpty())

	dynamic := make([]map[*dfanode]struct{}, k+1)
	dynamic[0] = map[*dfanode]struct{}{dfa.start: struct{}{}}

	for i := 1; i <= k; i++ {
		endpoint := false
		dynamic[i] = make(map[*dfanode]struct{})

		for node := range dynamic[i-1] {
			for _, next := range node.next {
				dynamic[i][next] = struct{}{}

				if next.endpoint {
					endpoint = true
				}
			}
		}

		if endpoint {
			if k%i == 0 {
				return true
			}
		}
	}

	return false
}

func (st *stack) Push(v any) {
	*st = append(*st, v)
}

func (st stack) Size() int {
	return len(st)
}

func (st stack) Top() any {
	return st[st.Size()-1]
}

func (st *stack) Pop() {
	(*st) = (*st)[0 : st.Size()-1]
}
