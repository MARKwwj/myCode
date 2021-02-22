package main

import "fmt"

func sum(values []int, resultChan chan int) {
	sum := 0
	for _, value := range values {
		sum += value
	}
	resultChan <- sum // 将计算结果发送到channel中
}
func main() {
	var values []int
	for i := 1; i < 999999999; i++ {
		values = append(values, i)
	}
	fmt.Println(values)
	resultChan := make(chan int, 2)
	go sum(values[:len(values)/2], resultChan)
	go sum(values[len(values)/2:], resultChan)
	sum1, sum2 := <-resultChan, <-resultChan // 接收结果
	fmt.Println("Result:", sum1, sum2, sum1+sum2)
}
