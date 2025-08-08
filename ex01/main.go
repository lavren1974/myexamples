package main

import "fmt"

func add(x, y int) int {
	return x + y
}

func main() {
	a := 10
	b := 20
	result := add(a, b)
	fmt.Printf("The sum of %d and %d is %d\n", a, b, result)
}