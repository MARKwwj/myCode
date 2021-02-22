package main

import "fmt"

func fib(n int) int {
	if n < 2 {
		return n
	}
	return fib(n-1) + fib(n-2)
}

func merge(A []int, m int, B []int, n int) {
	for kb,vb := range B {
		if m > 0 {
			for ka:=m+kb-1; ka>=0; ka-- {
				if vb < A[ka]  {
					A[ka+1] = A[ka]
					A[ka] = vb
				} else {
					A[ka+1] = vb
					break
				}
			}
		} else {
			A[kb] = vb
		}
	}
}
func main() {
	A := []int{1, 2, 3, 0, 0, 0}
	B := []int{2, 5, 6}
	merge(A, 3, B, 3)
	fmt.Println(A)
}
