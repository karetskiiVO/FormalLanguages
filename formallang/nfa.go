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
	nodes []*nfanode
}

type nfanode struct {
	next     map[rune]([](*nfanode))
	endpoint bool
}

// NFAFromRegExp - constructs new NFA with given regular expression
func NFAFromRegExp(reg *RegExp) *NFA {
	var res NFA

	res.abc = maps.Clone(reg.abc)

	return &res
}

func (nfa *NFA) newNode() *nfanode {
	res := nfanode{
		next:     make(map[rune][](*nfanode)),
		endpoint: false,
	}

	nfa.nodes = append(nfa.nodes, &res)
	return &res
}

// Dump - dumps NFA into png
func (nfa NFA) Dump(filename string) {
	g := graphviz.New()
	graph, err := g.Graph()
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
	for _, nodeptr := range nfa.nodes {
		graphnode, err := graph.CreateNode(fmt.Sprintf("%p", nodeptr))

		fromGRAFtoNFA[graphnode] = nodeptr
		fromNFAtoGRAF[nodeptr] = graphnode
		if err != nil {
			log.Fatal(err)
		}
	}

	for _, from := range nfa.nodes {
		for r, links := range from.next {
			for _, to := range links {
				edge, err := graph.CreateEdge(fmt.Sprintf("%p_%p", fromNFAtoGRAF[from], fromNFAtoGRAF[to]), fromNFAtoGRAF[from], fromNFAtoGRAF[to])
				edge.SetLabel(fmt.Sprintf("%c", r))

				if err != nil {
					log.Fatal(err)
				}
			}
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
