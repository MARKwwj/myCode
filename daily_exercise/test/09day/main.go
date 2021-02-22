package main

import (
	"fmt"
	"reflect"
)

func reflectType(v interface{}) {
	t := reflect.TypeOf(v)
	fmt.Println(t)
	fmt.Println(t.Name())
	fmt.Println(t.Kind())

}

func main() {
	//反射
	var str string
	reflectType(str)
	var r rune
	reflectType(r)

}
