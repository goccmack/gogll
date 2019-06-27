package main

import "fmt"

type A struct {
	I int
	J int
	K int
}

func main() {
	x, y := A{1, 2, 3}, A{1, 2, 4}
	fmt.Println(x == y)
	fmt.Println(x != y)
}
