package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io/ioutil"
)

func toMD5(mp3Byte []byte) string {
	m := md5.New()
	m.Write(mp3Byte)
	return hex.EncodeToString(m.Sum(nil))
}
func main() {
	mp3Byte1, _ := ioutil.ReadFile("D:\\desktop\\1.mp3")
	mp3Byte2, _ := ioutil.ReadFile("D:\\desktop\\2.mp3")
	s1 := toMD5(mp3Byte1)
	s2 := toMD5(mp3Byte2)
	fmt.Println(s1)
	fmt.Println(s2)
	if s1 == s2 {
		fmt.Println("true")
	} else {
		fmt.Println("false")
	}

}
