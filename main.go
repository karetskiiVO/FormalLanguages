package main

import (
	"bufio"
	"os"
	
	fl "github.com/karetskiiVO/FormalLanguages/formallang"
)

func main() {
	str, _ := bufio.NewReader(os.Stdin).ReadString('\n')

	reg, _ := fl.RegExpFromTokens(testConvert(str))
	aut := fl.NFAFromRegExp(reg)
	aut.Dump("./result0.png")
	aut.RemoveEmpty().Dump("./result1.png")
	daut := fl.DFAfromNFA(aut)
	daut.Dump("./result2.png")
}

func testConvert (str string) []fl.Token {
	res := make([]fl.Token, 0)

	special := map[rune]struct{} {
		'(': {},
		')': {},
		'+': {},
		'*': {},
		'1': {},
	}

	for _, r := range str {
		if r == ' ' || r == '\r' ||  r == '\n' {
			continue
		}

		_, ok := special[r]
		res = append(res, fl.Token{Symb: r, Servicable: ok});
	}

	return res
}
