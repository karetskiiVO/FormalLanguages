package main

import (
	"bufio"
	"fmt"
	fl "formallang"
	"os"
)

func main() {
	str, _ := bufio.NewReader(os.Stdin).ReadString('\n')

	reg, _ := fl.RegExpFromTokens(testConvert(str))
	fl.NFAFromRegExp(reg).Dump("./result.png")

	fmt.Println("aboba")
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
