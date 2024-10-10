package formallang

import (
	"fmt"
	"log"
	"maps"
	"sort"
	"strings"
	"unsafe"

	"github.com/goccy/go-graphviz"
	"github.com/goccy/go-graphviz/cgraph"
)

// DFA - imlement deterministic finite automaton with one letter transition
type DFA struct {
	abc   map[rune]struct{}
	nodes map[*dfanode]struct{}
	start *dfanode
}

type dfanode struct {
	next     map[rune]*dfanode
	linkscnt int
	endpoint bool
}

func (from *dfanode) link(r rune, to *dfanode) *dfanode {
	from.next[r] = to
	
	if from != to {
		to.linkscnt++
	}

	return from
}

func (from *dfanode) unlink(r rune, to *dfanode) *dfanode {
	if from != to {
		to.linkscnt--
	}

	delete(from.next, r)
	return from
}

func (dfa *DFA) newNode() *dfanode {
	res := dfanode{
		next:     make(map[rune]*dfanode),
		linkscnt: 0,
		endpoint: false,
	}

	dfa.nodes[&res] = struct{}{}
	return &res
}

// DFAfromNFA - constructs new NFA with DFA
func DFAfromNFA(nfa *NFA) *DFA {
	dfa := &DFA{
		abc:   maps.Clone(nfa.abc),
		nodes: make(map[*dfanode]struct{}),
	}

	SliceToString := func(sl []*nfanode) string {
		builder := &strings.Builder{}

		sort.Slice(sl, func(i, j int) bool {
			return fmt.Sprint(sl[i]) < fmt.Sprint(sl[j])
		})

		fmt.Fprintf(builder, "%v", len(sl))

		for _, ptr := range sl {
			fmt.Fprintf(builder, ",%p", ptr)
		}

		return builder.String()
	}
	StringToSlice := func(str string) []*nfanode {
		reader := strings.NewReader(str)
		size := 0
		fmt.Fscanf(reader, "%d", &size)
		res := make([]*nfanode, size)

		for i := 0; i < size; i++ {
			var ptr uintptr
			fmt.Fscanf(reader, ",%v", &ptr)
			res[i] = (*nfanode)(unsafe.Pointer(ptr))
		}

		return res
	}

	condition := append([]*nfanode{}, nfa.start)

	var tasks queue
	used := make(map[string]*dfanode)

	dfa.start = dfa.newNode()
	dfa.start.endpoint = nfa.start.endpoint
	tasks.Push(SliceToString(condition))
	used[SliceToString(condition)] = dfa.start

	for tasks.Size() > 0 {
		currCondString := tasks.Top().(string)
		tasks.Pop()

		currCond := StringToSlice(currCondString)
		dfafrom := used[currCondString]

		for r := range dfa.abc {
			nextCondSet := make(map[*nfanode]struct{})
			endpoint := false 

			for _, nfafrom := range currCond {
				for nfato := range nfafrom.next[r] {
					if nfato.endpoint {
						endpoint = true
					}

					nextCondSet[nfato] = struct{}{}
				}
			}

			if len(nextCondSet) == 0 {
				continue
			}

			nextCond := make([]*nfanode, 0, len(nextCondSet))
			for key := range nextCondSet {
				nextCond = append(nextCond, key)
			}

			nextCondString := SliceToString(nextCond)
			if _, ok := used[nextCondString]; !ok {

				node := dfa.newNode()
				node.endpoint = endpoint
				used[nextCondString] = node

				tasks.Push(nextCondString)
			}

			dfato := used[nextCondString]

			dfafrom.link(r, dfato)
		}
	}

	return dfa
}

// Dump - dumps NFA into png
func (dfa DFA) Dump(filename string) {
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

	fromNFAtoGRAF := make(map[*dfanode]*cgraph.Node)
	fromGRAFtoNFA := make(map[*cgraph.Node]*dfanode)
	for nodeptr := range dfa.nodes {
		graphnode, err := graph.CreateNode(fmt.Sprintf("%p", nodeptr))

		nodeShape := "circle"
		if nodeptr.endpoint {
			nodeShape = "doublecircle"
		}

		nodeLabel := ""
		if nodeptr == dfa.start {
			nodeLabel = "in"
		}

		graphnode.SetLabel(nodeLabel).SetShape(cgraph.Shape(nodeShape))

		fromGRAFtoNFA[graphnode] = nodeptr
		fromNFAtoGRAF[nodeptr] = graphnode
		if err != nil {
			log.Fatal(err)
		}
	}

	buf := make(map[struct{ from, to *dfanode }]([]rune))

	for from := range dfa.nodes {
		for r, to := range from.next {
			buf[struct{ from, to *dfanode }{from, to}] = append(buf[struct{ from, to *dfanode }{from, to}], r)
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
