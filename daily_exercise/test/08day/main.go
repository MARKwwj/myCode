package main

import "fmt"

type person struct {
	name  string
	age   int
	hobby string
}

func main() {
	p := newPerson("xx", "yy", 18)
	fmt.Printf("%p\n", p)
}

func newPerson(name, hobby string, age int) *person {
	p := &person{
		name:  name,
		age:   age,
		hobby: hobby,
	}
	return p
}
