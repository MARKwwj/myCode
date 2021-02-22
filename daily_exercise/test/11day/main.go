package main

import (
	"fmt"
	"math/rand"
)

//var wg sync.WaitGroup

//func hello(i int) {
//	defer wg.Done() // goroutine结束就登记-1
//	fmt.Println("Hello Goroutine!", i)
//}
func main() {
	//
	//for i := 0; i < 10; i++ {
	//	wg.Add(1) // 启动一个goroutine就登记+1
	//	go hello(i)
	//}
	////wg.Wait() // 等待所有登记的goroutine都结束
	//err:=os.Rename("D:\\desktop\\r2", "D:\\desktop\\ssss")
	//fmt.Println(err)
	for  {
		//fmt.Println(rand.Intn(9))
		//fmt.Println(rand.Intn(87000) + 12900)
		fmt.Println((float32(rand.Intn(9)) / 10) + float32(8.8))
	}
}
