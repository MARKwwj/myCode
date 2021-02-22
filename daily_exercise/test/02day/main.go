package main

import "fmt"

func main() {
	//map
	mp1 := make(map[string]int, 5)
	mp1["id"] = 10
	mp1["i"] = 9
	mp1["cd"] = 8
	fmt.Println(mp1)
	//map声明时赋值
	mp2 := map[string]string{
		"name":        "zhangsan",
		"add":         "shanghai",
		"phonenumber": "18875695214",
	}
	fmt.Println(mp2)
	delete(mp2, "name")
	for k, v := range mp2 {
		fmt.Println(k, v)
	}
	value, ok := mp2["ss"]
	if ok {
		fmt.Println(ok, value)
	} else {
		fmt.Println(ok, value)
	}

	mp3 := map[string]int{
		"a": 1,
		"b": 2,
		"c": 3,
		"d": 4,
	}
	for k, v := range mp3 {
		fmt.Println(k, v)

	}
	var mp4 map[string]int
	mp4 = make(map[string]int)
	mp4["a"] = 1
	mp4["b"] = 2
	mp4["c"] = 3
	mp4["d"] = 4
	for k, v := range mp4 {
		fmt.Println(k, v)
	}
	//元素为map 类型的切片
	var mapSlice = make([]map[string]int, 10)
	for index,value :=range mapSlice{
		fmt.Println(index,value)
	}
	mapSlice[0]=make(map[string]int)
	mapSlice[0]["id"]=1
	mapSlice[0]["name"]=1
	mapSlice[0]["addr"]=1
	mapSlice[0]["time"]=1
	for index,value :=range mapSlice{
		fmt.Println(index,value)
	}

}
