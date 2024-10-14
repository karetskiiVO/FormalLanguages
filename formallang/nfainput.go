package formallang

import "maps"

// NFANodeInput - struct with description of every node
type NFANodeInput struct {
	id   string
	next []struct {
		r  rune
		id string
	}
	start    bool
	endpoint bool
}

// NFAfromInput - creates NFA from given array of nodes descriptions
func NFAfromInput(abc map[rune]struct{}, input []NFANodeInput) *NFA {
	nfa := &NFA{
		abc: maps.Clone(abc),
	}

	idToNodes := make(map[string]*nfanode)
	for _, descr := range input {
		if _, ok := idToNodes[descr.id]; !ok {
			idToNodes[descr.id] = nfa.newNode()
		}

		from := idToNodes[descr.id]
		from.endpoint = descr.endpoint

		if descr.start {
			nfa.start = from
		}

		for _, nextnode := range descr.next {
			if _, ok := idToNodes[nextnode.id]; !ok {
				idToNodes[nextnode.id] = nfa.newNode()
			}

			to := idToNodes[nextnode.id]

			from.link(nextnode.r, to)
		}
	}

	return nfa
}
