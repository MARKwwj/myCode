package main

import "fmt"

func main() {

	a := 1
	var d = 3
	fmt.Print(d)
	var b int = 3
	fmt.Print(b)
	fmt.Print(a)
	fmt.Println()
	s := '1'
	s1 := "sss"

	fmt.Printf("%T \n", s)
	fmt.Printf("%T \n", s1)

	//方式一
	var arr = [3]int{1, 2, 3}
	for _, r := range arr {
		fmt.Println(r)
	}
	//方式二
	var arr2 = [...]int{1, 2, 3}
	for _, r := range arr2 {
		fmt.Println(r)
	}
	//方式三
	arr3 := [...]string{"花花", "世界", "简介", "搜索"}
	for _, r := range arr3 {
		fmt.Println(r)
	}

	//多维数组
	arr4 := [...][5]string{
		{"花花", "世界", "简介", "搜索"},
		{"花花", "世界", "简介", "搜索"},
	}
	fmt.Println(arr4[1][1])

	//切片
	p1 := []int{1, 3, 4, 5, 6, 7}
	fmt.Println(p1[1:3])
	// cap 容量
	fmt.Println(len(p1), cap(p1))

	arr9 := [...]int{1, 2, 3, 4, 5, 6, 7}
	p2 := arr9[1:3:3]
	fmt.Println(len(p2), cap(p2))

	//使用make() 函数 构造切片
	p3 := make([]int, 5, 8)
	fmt.Println(p3)
	fmt.Println(len(p3), cap(p3))
	p4 := make([]string, 5, 6)
	p4 = []string{"试试", "搜索", "手术", "地方"}
	fmt.Println(p4[1:3:3])

	//切片的复制拷贝

	// 切片的底层是一个数组  切片是引用类型 数组是值类型
	k1 := make([]int, 5)
	k2 := k1
	k2[0] = 100
	fmt.Println(k1)
	fmt.Println(k2)
	//遍历 切片 方式一
	for _, r := range k1 {
		fmt.Println(r)
	}
	fmt.Println("方式二")
	//遍历 切片 方式二
	for i := 0; i < len(k2); i++ {
		fmt.Println(i, k2[i])

	}
	//append 为切片添加元素
	l1 := []int{33, 44, 55}
	k2 = append(k2, l1...)
	fmt.Println(k2)
	// 使用 var 声明切时 不需要再初始化
	var l2 []int
	l2 = append(l2, 1,2,3,4)
	fmt.Println(l2)


}
