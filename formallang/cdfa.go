package formallang

import (
	"fmt"
	"log"
	"maps"
	"strings"

	"github.com/goccy/go-graphviz"
	"github.com/goccy/go-graphviz/cgraph"
)

// CDFA - imlement complete deterministic finite state automaton with
type CDFA struct {
	abc   map[rune]struct{}
	nodes map[*dfanode]struct{}
	start *dfanode
	stock *dfanode
}

// CDFAfromDFA - constructs new CDFA from DFA
func CDFAfromDFA(dfa *DFA) *CDFA {
	cdfa := &CDFA{
		abc:   maps.Clone(dfa.abc),
		nodes: make(map[*dfanode]struct{}),
	}

	DFAtoCDFA := make(map[*dfanode]*dfanode)
	cdfa.stock = cdfa.newNode()
	for r := range cdfa.abc {
		cdfa.stock.link(r, cdfa.stock)
	}

	for dfanode := range dfa.nodes {
		cdfanode := cdfa.newNode()
		cdfanode.endpoint = dfanode.endpoint
		DFAtoCDFA[dfanode] = cdfanode
	}

	cdfa.start = DFAtoCDFA[dfa.start]

	for dfafrom, cdfafrom := range DFAtoCDFA {
		for r, dfato := range dfafrom.next {
			cdfato := DFAtoCDFA[dfato]

			cdfafrom.link(r, cdfato)
		}

		for r := range dfa.abc {
			if _, ok := cdfafrom.next[r]; ok {
				continue
			}

			cdfafrom.link(r, cdfa.stock)
		}
	}

	return cdfa
}

// DFAfromCDFA - constructs new DFA from CDFA
func DFAfromCDFA(cdfa *CDFA) *DFA {
	dfa := &DFA{
		abc:   maps.Clone(cdfa.abc),
		nodes: make(map[*dfanode]struct{}),
	}

	CDFAtoDFA := make(map[*dfanode]*dfanode)

	for cdfanode := range cdfa.nodes {
		if cdfanode == cdfa.stock {
			continue
		}

		dfanode := dfa.newNode()
		dfanode.endpoint = cdfanode.endpoint
		CDFAtoDFA[cdfanode] = dfanode
	}

	dfa.start = CDFAtoDFA[cdfa.start]

	for cdfafrom, dfafrom := range CDFAtoDFA {
		for r, cdfato := range cdfafrom.next {
			if cdfato == cdfa.stock {
				continue
			}

			dfato := CDFAtoDFA[cdfato]
			dfafrom.link(r, dfato)
		}
	}

	return dfa
}

// Minimise constructs mdfa
func (cdfa CDFA) Minimise() *CDFA {
	nodeClasses := make(map[*dfanode]string)
	classesSet := make(map[string]int)

	for node := range cdfa.nodes {
		class := "0"
		if node.endpoint {
			class = "1"
		}

		nodeClasses[node] = class
		classesSet[class] = 1
	}

	alph := make([]rune, 0, len(cdfa.abc))
	for r := range cdfa.abc {
		alph = append(alph, r)
	}

	for {
		bufNodeClasses := make(map[*dfanode]string)
		bufClassesSet := make(map[string]int)

		cnt := 0
		for from, fromclass := range nodeClasses {
			newClassBuilder := &strings.Builder{}
			
			newClassBuilder.WriteString(fromclass)
			for _, r := range alph {
				to := from.next[r]
				newClassBuilder.WriteRune(',')
				newClassBuilder.WriteString(nodeClasses[to])
			}

			newClass := newClassBuilder.String()

			var classid int
			if id, ok := bufClassesSet[newClass]; !ok {
				bufClassesSet[newClass] = cnt
				classid = cnt
				cnt++
			} else {
				classid = id
			}

			bufNodeClasses[from] = fmt.Sprint(classid)
		}

		if (len(bufClassesSet) == len(classesSet)) {
			classesSet = bufClassesSet
			nodeClasses = bufNodeClasses
			break
		}

		classesSet = bufClassesSet
		nodeClasses = bufNodeClasses
	}

	classesSet = make(map[string]int)
	for _, class := range nodeClasses {
		classesSet[class] = 1
	}

	mcdfa := &CDFA{
		abc:   maps.Clone(cdfa.abc),
		nodes: make(map[*dfanode]struct{}),
	}

	classesToMCDFANodes := make(map[string]*dfanode)
	for class := range classesSet {
		classesToMCDFANodes[class] = mcdfa.newNode()
	}

	mcdfa.start = classesToMCDFANodes[nodeClasses[cdfa.start]]
	mcdfa.stock = classesToMCDFANodes[nodeClasses[cdfa.stock]]

	for cdfafrom, fromclass := range nodeClasses {
		mcdfafrom := classesToMCDFANodes[fromclass]
		
		if cdfafrom.endpoint {
			mcdfafrom.endpoint = true
		}

		for r, cdfato := range cdfafrom.next {
			toclass := nodeClasses[cdfato]
			mcdfato := classesToMCDFANodes[toclass]

			mcdfafrom.link(r, mcdfato)
		}
	}

	return mcdfa
}

func (cdfa *CDFA) newNode() *dfanode {
	res := dfanode{
		next:     make(map[rune]*dfanode),
		linkscnt: 0,
		endpoint: false,
	}

	cdfa.nodes[&res] = struct{}{}
	return &res
}

// Dump - dumps CDFA into png
func (cdfa CDFA) Dump(filename string) {
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
	for nodeptr := range cdfa.nodes {
		graphnode, err := graph.CreateNode(fmt.Sprintf("%p", nodeptr))

		nodeShape := "circle"
		if nodeptr.endpoint {
			nodeShape = "doublecircle"
		}

		nodeLabel := ""
		if nodeptr == cdfa.start {
			nodeLabel = "in"
		}
		if nodeptr == cdfa.stock {
			nodeLabel = "stok"
		}

		graphnode.SetLabel(nodeLabel).SetShape(cgraph.Shape(nodeShape))

		fromGRAFtoNFA[graphnode] = nodeptr
		fromNFAtoGRAF[nodeptr] = graphnode
		if err != nil {
			log.Fatal(err)
		}
	}

	buf := make(map[struct{ from, to *dfanode }]([]rune))

	for from := range cdfa.nodes {
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
