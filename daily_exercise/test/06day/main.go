package main

import "fmt"

type calculaction func(int, int) int
type myint int  //基于int的新类型
type myint2 = myint //类型别名

func main() {
	//var f calculaction
	//f = add
	//res := f(1, 2)
	//fmt.Println(res)
	//a := []int{1, 2, 3}
	//for _, v := range a {
	//	value := v
	//	fmt.Println(&value)
	//}
	//
	//	sn1 := struct {
	//		age  int
	//		name string
	//	}{age: 11, name: "qq"}
	//	sn2 := struct {
	//		age  int
	//		name string
	//	}{age: 11, name: "qq"}
	//
	//	if sn1 == sn2 {
	//		fmt.Println("sn1 == sn2")
	//	}

	sm1 := struct {
		age int
		m   int
	}{age: 11, m: 1}
	sm2 := struct {
		age int
		m   int
	}{age: 11, m: 1}

	if sm1 == sm2 {
		fmt.Println("sm1 == sm2")
	}
}

func add(a, b int) (res int) {
	return a + b
}
