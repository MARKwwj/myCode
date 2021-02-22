package main

import "fmt"

func main() {
	//res := m(1, 2)
	//fmt.Println(res)
	//a := 100
	//b := 200
	//fmt.Printf("a:%v ", a)
	//fmt.Printf("&a:%v\n", &a)
	//
	//fmt.Printf("b:%v ", b)
	//fmt.Printf("&b:%v\n", &b)
	//
	//changeValue(a, b)
	//
	//fmt.Printf("交换后的值：a:%v,b:%v \n \n ", a, b)\
	a := 100
	b := 200
	fmt.Println(&a)
	fmt.Println(&b)

	c(a,b)
}

func m(a int, b int) string {
	fmt.Println(a, b)
	return "ok"

}

func changeValue(a, b int) {
	var temp int
	temp = a
	fmt.Printf("temp:%v ", temp)
	fmt.Printf("&temp:%v\n", &temp)
	a = b
	fmt.Printf("a:%v ", a)
	fmt.Printf("&a:%v\n", &a)
	b = temp
	fmt.Printf("b:%v ", b)
	fmt.Printf("&b:%v\n", &b)
}
