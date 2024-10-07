package main

import (
	"fmt"
	fl "formallang"
)

func main() {
	fmt.Println(fl.Test(
		[]fl.Token{
			{Symb: 'a', Servicable: false},
			{Symb: 'a', Servicable: false},
			{Symb: '+', Servicable: false},
		},
	))
}
