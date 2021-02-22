package main

import (
	"fmt"
)

//func CheckPermutation(s1 string, s2 string) bool {
//	f := func(s string) []byte {
//		byteS1 := []byte(s)
//		for i := 0; i < len(byteS1)-1; i++ {
//			for j := i; j < len(byteS1); j++ {
//				if byteS1[i] > byteS1[j] {
//					temp := byteS1[j]
//					byteS1[j] = byteS1[i]
//					byteS1[i] = temp
//				}
//			}
//		}
//		return byteS1
//	}
//	ss := string(f(s1))
//	ss2 := string(f(s2))
//	if ss == ss2 {
//		return true
//	} else {
//		return false
//	}
//}
//func CheckPermutation(s1 string, s2 string) bool {
//	if len(s1) != len(s2) {
//		return false
//	}
//	tmp1 := strings.Split(s1, "")
//	tmp2 := strings.Split(s2, "")
//	sort.Strings(tmp1)
//	res1 := strings.Join(tmp1, "")
//	sort.Strings(tmp2)
//	res2 := strings.Join(tmp2, "")
//	return res1 == res2
//}
//魔术索引。 在数组A[0...n-1]中，有所谓的魔术索引，满足条件A[i] = i。给定一个有序整数数组，
//编写一种方法找出魔术索引，若有的话，在数组A中找出一个魔术索引，如果没有，则返回-1。
//若有多个魔术索引，返回索引值最小的一个。
func findMagicIndex(nums []int) int {
	for k, v := range nums {
		if v == k {
			return k
		}
	}
	return -1
}
func main() {
	n := findMagicIndex([]int{1, 3, 4, 3, 2, 5})
	fmt.Println(n)
}
