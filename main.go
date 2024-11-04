package main

import (
	"fmt"

	fl "github.com/karetskiiVO/FormalLanguages/formallang"
)

func main() {
	var pol string
	var r rune
	var k int

	fmt.Scanf("%s %c %d", &pol, &r, &k)

	if fl.Solve7(pol, r, k) {
		fmt.Println("YES")
	} else {
		fmt.Println("NO")
	}
}
