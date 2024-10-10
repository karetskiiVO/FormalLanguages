package formallang

import (
	"fmt"
	"log"
	"maps"

	"github.com/goccy/go-graphviz"
	"github.com/goccy/go-graphviz/cgraph"
)

const (
	// EmptyRune - is empty transition move key
	EmptyRune = rune('$')
)

// NFA - imlement nondeterministic finite automaton with one letter transition
type NFA struct {
	abc   map[rune]struct{}
	nodes map[*nfanode]struct{}
	start *nfanode
}

type nfanode struct {
	next     map[rune]map[*nfanode]struct{}
	linkscnt int
	endpoint bool
}

func (from *nfanode) link(r rune, to *nfanode) *nfanode {
	if from.next[r] == nil {
		from.next[r] = map[*nfanode]struct{}{}
	}

	if from != to {
		to.linkscnt++
	}

	from.next[r][to] = struct{}{}
	return from
}

func (from *nfanode) unlink(r rune, to *nfanode) *nfanode {
	if from.next[r] == nil {
		return from
	}

	if from != to {
		to.linkscnt--
	}

	delete(from.next[r], to)
	return from
}

// NFAFromRegExp - constructs new NFA with given regular expression
func NFAFromRegExp(reg *RegExp) *NFA {
	res := &NFA{
		nodes: make(map[*nfanode]struct{}),
	}

	res.abc = maps.Clone(reg.abc)
	begin, end := res.newNode(), res.newNode()
	res.start = begin
	res.start.linkscnt = 1
	end.endpoint = true

	reg.tree.ToSubNFA(res, begin, end)

	return res
}

// RemoveEmpty - removes emty links
func (nfa *NFA) RemoveEmpty() *NFA {
	for from := range nfa.nodes {
		if _, ok := from.next[EmptyRune][from]; ok {
			delete(from.next[EmptyRune], from)
		}
	}

	for from := range nfa.nodes {
		emptyReleased := false

		for !emptyReleased {
			emptyReleased = true
			for to := range from.next[EmptyRune] {
				if to.endpoint {
					from.endpoint = true
				}

				from.unlink(EmptyRune, to)
				for r, tonext := range to.next {
					if r == EmptyRune {
						emptyReleased = false
					}

					for node := range tonext {
						from.link(r, node)
					}
				}
			}
		}
	}

	nfa.removeNoLinks()

	return nfa
}

func (nfa *NFA) newNode() *nfanode {
	res := nfanode{
		next:     make(map[rune]map[*nfanode]struct{}),
		linkscnt: 0,
		endpoint: false,
	}

	nfa.nodes[&res] = struct{}{}
	return &res
}

func (nfa *NFA) deleteNode(node *nfanode) {
	delete(nfa.nodes, node)
}

func (nfa *NFA) removeNoLinks() {
	for node := range nfa.nodes {
		if node.linkscnt == 0 {
			nfa.deleteNode(node)
		}
	}
}

// Dump - dumps NFA into png
func (nfa NFA) Dump(filename string) {
	g := graphviz.New()
	graph, err := g.Graph(graphviz.StrictDirected)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := graph.Close(); err != nil {
			log.Fatal(err)
		}
		g.Close()
	}()

	fromNFAtoGRAF := make(map[*nfanode]*cgraph.Node)
	fromGRAFtoNFA := make(map[*cgraph.Node]*nfanode)
	for nodeptr := range nfa.nodes {
		graphnode, err := graph.CreateNode(fmt.Sprintf("%p", nodeptr))

		nodeShape := "circle"
		if nodeptr.endpoint {
			nodeShape = "doublecircle"
		}

		nodeLabel := ""
		if nodeptr == nfa.start {
			nodeLabel = "in"
		}

		graphnode.SetLabel(nodeLabel).SetShape(cgraph.Shape(nodeShape))

		fromGRAFtoNFA[graphnode] = nodeptr
		fromNFAtoGRAF[nodeptr] = graphnode
		if err != nil {
			log.Fatal(err)
		}
	}

	buf := make(map[struct{ from, to *nfanode }]([]rune))

	for from := range nfa.nodes {
		for r, links := range from.next {
			for to := range links {
				buf[struct{ from, to *nfanode }{from, to}] = append(buf[struct{ from, to *nfanode }{from, to}], r)
			}
		}
	}

	for pair, runes := range buf {
		edge, err := graph.CreateEdge(fmt.Sprintf("%p_%p", fromNFAtoGRAF[pair.from], fromNFAtoGRAF[pair.to]), fromNFAtoGRAF[pair.from], fromNFAtoGRAF[pair.to])

		lable := fmt.Sprintf("%c", runes[0])
		for _, r := range runes[1:] {
			lable = fmt.Sprintf("%s,%c", lable, r)
		}

		edge.SetLabel(lable)

		if err != nil {
			log.Fatal(err)
		}
	}

	// 1. write encoded PNG data to buffer
	//var buf bytes.Buffer
	//if err := g.Render(graph, graphviz.PNG, &buf); err != nil {
	//	log.Fatal(err)
	//}
	//
	//// 2. get as image.Image instance
	//image, err := g.RenderImage(graph)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//// 3. write to file directly
	if err := g.RenderFilename(graph, graphviz.PNG, filename); err != nil {
		log.Fatal(err)
	}
}
