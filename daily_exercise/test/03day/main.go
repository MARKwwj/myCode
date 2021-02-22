package main

import "fmt"

func main() {
	//写一个程序，统计一个字符串中每个单词出现的次数。
	//比如：”how do you do”中how=1 do=2 you=1。
	//var str string
	//str = "where there is a will there is a way"
	//strSlice := strings.Split(str, " ")
	//fmt.Println(strSlice)
	//countMap := make(map[string]int, 10)
	//for _, v := range strSlice {
	//	countMap[v] += 1
	//}
	//fmt.Println(countMap)
		type Map map[string][]int
		m := make(Map)
		s := []int{1, 2}
		s = append(s, 3)
		fmt.Printf("%+v\n", s)
		m["q1mi"] = s
		s = append(s[:1], s[2:]...)
		fmt.Printf("%+v\n", s)
		fmt.Printf("%+v\n", m["q1mi"])
}
