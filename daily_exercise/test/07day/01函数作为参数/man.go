package main

import "fmt"

func main() {
	res := calc(1, 2, add)
	fmt.Println(res)
}
func add(a, b int) int {
	return a + b
}
func calc(a, b int, op func(int, int) int) int {
	return op(a, b)
}
