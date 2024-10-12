package main

import (
	"fmt"
	"os"

	fl "github.com/karetskiiVO/FormalLanguages/formallang"
)

func main() {
	makeTests()
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

func makeTests () {
	tests := []string{
		"a", "a+b", "ab", "a+1", "(a+1)*", 
		"(a(a+1) + b)*", "(a+b)*a(b+a)", "a((ba)*a(ab)* + a)*", 
		"(a(ab + ba)*b(a + ba)*)(a(ab + ba)*b(a + ba)*)*",
	}
	
	for idx, testregexp := range tests {
		makeTest(testregexp, "test"+fmt.Sprint(idx))
	}
}

func makeTest (regExp string, foldername string) {
	os.Mkdir("./test/"+foldername, os.ModeDir)

	reg, _ := fl.RegExpFromTokens(testConvert(regExp))
	aut := fl.NFAFromRegExp(reg)
	aut.Dump("./test/"+foldername+"/0_nfa.png")
	aut.RemoveEmpty().Dump("./test/"+foldername+"/1_nfa_without_empty.png")
	daut := fl.DFAfromNFA(aut)
	daut.Dump("./test/"+foldername+"/2_dfa.png")
	cdaut := fl.CDFAfromDFA(daut)
	cdaut.Dump("./test/"+foldername+"/3_cdfa.png")
	mfdaut := cdaut.Minimise()
	mfdaut.Dump("./test/"+foldername+"/4_mcdfa.png")
}