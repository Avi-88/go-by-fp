package main

import (
	"errors"
	"fmt"
)

// Day 2: Pointers and Memory
//
// Read EXERCISE.md before starting.
// Implement all three parts in this file (or split across files if you prefer).

type Node struct {
	value int
	next *Node
}

type List struct {
	head *Node
	size int  
}

func NewList () *List {
	return &List{}
}

func (l *List) Push (value int) {
	var node Node = Node{value: value, next: l.head}
	l.head = &node
	l.size++
}

func (l *List) Pop () (value int, err error) {
	if l.size == 0 {
		return 0,errors.New("The list is empty")
	}
	var nodeVal int = l.head.value
	l.head = l.head.next
	l.size--

	return nodeVal, nil
}

func (l *List) Append (value int) {
	if l.head == nil {
		var nNode Node = Node{value: value, next: nil}
		l.head = &nNode
		l.size++
		return
	}
	var tail *Node = l.head
	for tail.next != nil {
		tail = tail.next
	}
	var nNode Node = Node{value: value, next: nil}
	tail.next = &nNode
	l.size++
}

func (l *List) Len () int {
	return l.size
}

func (l *List) Print() {
	if l.head == nil {
		fmt.Print("[]")
		return
	}
	var tNode *Node = l.head
	fmt.Print("[")
	for tNode.next != nil {
		fmt.Printf("%d->", tNode.value)
		tNode = tNode.next
	}
	fmt.Printf("%d]\n",tNode.value)
}

func ShowAddresses() {
	var x int = 42
	var p *int = &x
	fmt.Printf("x value: %d\n", x)
	fmt.Printf("x address: %d\n", &x)
	fmt.Printf("p value: %d\n", p)
	fmt.Printf("*p value: %d\n", *p)

	*p = 100
	fmt.Printf("after *p = %d, x = %d\n", *p, x)
}

func doubleByValue(n int) int {
	n = n*2
	return n
}

func doubleByPointer(n *int) int {
	*n = (*n)*2
	return *n
}

func doubleSlice(s []int) []int {
	num := 0
	for num < len(s) {
		s[num] = s[num]*2
		num++
	}
	return s
}


func main() {
	list := NewList()
	list.Append(1)
	list.Append(2)
	list.Append(3)
	list.Push(0)
	list.Print()         // [0 -> 1 -> 2 -> 3]
	fmt.Println(list.Len()) // 4
	v, _ := list.Pop()
	fmt.Println(v)       // 0
	list.Print()         // [1 -> 2 -> 3]

	ShowAddresses()
	n := 5
	dv := doubleByValue(n)
	fmt.Printf("Before: %d and After: %d\n", n, dv)

	dp := doubleByPointer(&n)
	fmt.Printf("Before: %d and After: %d\n", n, dp)

	s := []int{1,2,3,4}
	fmt.Print("Before\n")
	for _,v := range s {
		fmt.Printf("%d\n", v)
	}

	ds := doubleSlice(s)
	fmt.Print("After\n")
	for _,v := range s {
		fmt.Printf("%d\n", v)
	}
	fmt.Println("-----------")
	for _,v := range ds {
		fmt.Printf("%d\n", v)
	}

}
