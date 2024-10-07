package formallang

import (
	"fmt"
)

type regExpNode interface {
	ToString(priority int) string
	Priority() int
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

type regExpNodeEmptyRune struct {}
func (regExpNodeEmptyRune) Priority() int         { return hightPriority }
func (r regExpNodeEmptyRune) ToString(int) string { return "1" }

type regExpNodeRune struct {
	r rune
}
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
