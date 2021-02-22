package main

import "fmt"

var (
	prices []int
	fee    int
)

func main() {
	//输入: prices = [1, 3, 2, 8, 4, 9], fee = 2
	//输出: 8
	//解释: 能够达到的最大利润:
	//在此处买入prices[0] = 1
	//在此处卖出 prices[3] = 8
	//在此处买入 prices[4] = 4
	//在此处卖出 prices[5] = 9
	//总利润:   ((8 - 1) - 2) + ((9 - 4) - 2) = 8
	prices = []int{1, 3, 2, 8, 4, 9}
	fee = 2
	for i := 0; i < len(prices); i++ {
		index1 := i
		for j := i; j < len(prices); j++ {
			index2 := j
			fmt.Print(index1, " ", index2," /")
		}
		fmt.Println()

	}
}
