package formallang

import (
	"fmt"
	"strings"
)

type regExpNode interface {
	ToString(priority int) string
	Priority() int
	ToSubNFA(nfa *NFA, begin, end *nfanode)
}

const (
	lowPriority = iota
	addPriority
	mulPriority
	cliniPriority
	runePriority
	emptyRunepriority
	hightPriority
)

type regExpNodeEmptyRune struct{}

func (regExpNodeEmptyRune) Priority() int       { return hightPriority }
func (regExpNodeEmptyRune) ToString(int) string { return "1" }
func (regExpNodeEmptyRune) ToSubNFA(nfa *NFA, begin, end *nfanode) {
	begin.link(EmptyRune, end)
}

type regExpNodeRune struct {
	r rune
}

func (regExpNodeRune) Priority() int         { return runePriority }
func (r regExpNodeRune) ToString(int) string { return fmt.Sprintf("%c", r.r) }
func (r regExpNodeRune) ToSubNFA(nfa *NFA, begin, end *nfanode) {
	begin.link(r.r, end)
}

type regExpNodeAdd struct {
	Next []regExpNode
}

func (regExpNodeAdd) Priority() int { return addPriority }
func (add regExpNodeAdd) ToString(priority int) string {
	prior := add.Priority()

	var builder strings.Builder

	if prior < priority {
		builder.WriteRune('(')
	}

	builder.WriteString(add.Next[0].ToString(prior))

	for _, next := range add.Next[1:] {
		builder.WriteString(" + ")
		builder.WriteString(next.ToString(prior))
	}

	if prior < priority {
		builder.WriteRune(')')
	}

	return builder.String()
}
func (add regExpNodeAdd) ToSubNFA(nfa *NFA, begin, end *nfanode) {
	for _, regexprnode := range add.Next {
		regexprnode.ToSubNFA(nfa, begin, end)
	}
}

type regExpNodeMul struct {
	Next []regExpNode
}

func (regExpNodeMul) Priority() int { return mulPriority }
func (mul regExpNodeMul) ToString(priority int) string {
	prior := mul.Priority()

	var builder strings.Builder

	if prior < priority {
		builder.WriteRune('(')
	}

	builder.WriteString(mul.Next[0].ToString(prior))

	for _, next := range mul.Next[1:] {
		builder.WriteString(next.ToString(prior))
	}

	if prior < priority {
		builder.WriteRune(')')
	}

	return builder.String()
}
func (mul regExpNodeMul) ToSubNFA(nfa *NFA, begin, end *nfanode) {
	nodes := make([]*nfanode, len(mul.Next)+1)
	nodes[0] = begin
	nodes[len(nodes)-1] = end

	for i := 1; i < len(mul.Next); i++ {
		nodes[i] = nfa.newNode()
	}

	for i, regexprnode := range mul.Next {
		regexprnode.ToSubNFA(nfa, nodes[i], nodes[i+1])
	}
}

type regExpNodeClini struct {
	Next regExpNode
}

func (regExpNodeClini) Priority() int { return cliniPriority }
func (clini regExpNodeClini) ToString(priority int) string {
	prior := clini.Priority()
	return fmt.Sprintf("%v*", clini.Next.ToString(prior))
}
func (clini regExpNodeClini) ToSubNFA(nfa *NFA, begin, end *nfanode) {
	clini.Next.ToSubNFA(nfa, begin, begin)
	begin.link(EmptyRune, end)
}