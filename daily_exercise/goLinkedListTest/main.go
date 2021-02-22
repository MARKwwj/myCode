package main

//空接口可存储任意类型
type Item interface {
}

//一个节点，除了自身的数据之外，还须指向下一个节点，尾部节点指向nil
type LinkNode struct {
	Payload Item //Payload 为任意数据类型
	Next    *LinkNode
}
type LinkNoder interface {
	Add(payload Item)
	Delete(index int) Item
}

func (head *LinkNode) Add(payload Item) {
	//采用从尾部插入的方式 给链表添加元素
	point := head
	for point.Next != nil {
		point = point.Next
	}
	newNode := LinkNode{Payload: payload, Next: nil}
	point.Next = &newNode
}

func main() {

}
