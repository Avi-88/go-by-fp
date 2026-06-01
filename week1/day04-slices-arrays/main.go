package main

import (
	"fmt"
	"errors"
)
// Day 4: Slices vs Arrays (Internals)
//
// Read EXERCISE.md before starting.
// Implement all 3 parts in this file.


// applies f to every element, returns new slice
func Map(s []int, f func(int) int) []int {
	newS := []int{}
	for _, i := range s {
		newI := f(i)
		newS = append(newS, newI)
	}
	return newS
}

// returns new slice with only elements where f returns true
func Filter(s []int, f func(int) bool) []int {
	newS := []int{}
	for _, i := range s {
		isTrue := f(i)
		if isTrue {
			newS = append(newS, i)
		}
	}
	return newS
}

// folds the slice into a single value
func Reduce(s []int, initial int, f func(int, int) int) int {
	init := initial

	for _,i := range s {
		init = f(init, i)
	}

	return init
}


// returns true if val is in s
func Contains(s []int, val int) bool {
	for _,i := range s {
		if i == val {
			return true
		}
	}
	return false
}

// returns a new reversed slice (do NOT modify the original)
func Reverse(s []int) []int {
	n := len(s)
	i := n-1
	newS := []int{}
	for {
		if i < 0 {
			break
		}
		newS = append(newS, s[i])
		i--
	}
	return newS
}

// returns a new slice with duplicates removed, preserving order
func Unique(s []int) []int {
	m := make(map[int]int)
	newS := []int{}
	for _,i := range s {
		_, ok := m[i]
		if !ok {
			newS = append(newS, i)
			m[i]++
		}
	}
	return newS
}

type DynamicArray struct {
	data []int
	length int
	capacity int
}

func NewDynamicArray() *DynamicArray {
	dArray := DynamicArray{
		data: make([]int, 0, 2),
		capacity: 2,
		length: 0,
	}
	return &dArray
}

func (d *DynamicArray) Push(val int) {
	if d.length == d.capacity {
		d.capacity = 2*d.capacity
		newDArray := make([]int, d.length, d.capacity)
		copy(newDArray, d.data)
		d.data = newDArray
		fmt.Printf("Capacity doubled to %v\n", d.capacity)
	}
	d.data = append(d.data, val)
	d.length++
}

func (d *DynamicArray) Get(i int) (int, error) {
	if i < 0 || i >= d.length {
		return 0, errors.New("Out of bounds index accessed")
	}
	return d.data[i], nil
}

func (d *DynamicArray) Len() int {
	return d.length
}

func (d *DynamicArray) Cap() int {
	return d.capacity
}

func (d *DynamicArray) Print() {
	for _,i := range d.data {
		fmt.Println(i)
	}
}


func main() {

	// //part 1
	// var arr [5]int = [5]int{1,2,3,4,5}

	// s := arr[1:4]
	// s[0] = 100
	// fmt.Println(s[0])
	// fmt.Println(arr[1])
	// s = append(s, 200)
	// s = append(s, 201)
	// s = append(s, 202)
	// s[0] = 130
	// fmt.Println(s[0])
	// fmt.Println(arr[1])


	// //part 2
	// nums := []int{1, 2, 3, 4, 5, 2, 3, 1}

	// doubled  := Map(nums, func(n int) int { return n * 2 })
	// evens    := Filter(nums, func(n int) bool { return n%2 == 0 })
	// sum      := Reduce(nums, 0, func(acc, n int) int { return acc + n })
	// reversed := Reverse(nums)
	// unique   := Unique(nums)

	// fmt.Println(doubled)   // [2 4 6 8 10 4 6 2]
	// fmt.Println(evens)     // [2 4 2]
	// fmt.Println(sum)       // 21
	// fmt.Println(reversed)  // [1 3 2 5 4 3 2 1]
	// fmt.Println(unique)    // [1 2 3 4 5]

	//part 3
	dArray := NewDynamicArray()

	dArray.Push(1)
	dArray.Push(2)
	dArray.Push(3)
	dArray.Push(4)
	dArray.Push(5)
	dArray.Push(6)

	idx0, err := dArray.Get(0)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Get 0 index", idx0)
	}
	

	idx4, err := dArray.Get(4)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Get 4 index", idx4)
	}
	

	idx6, err := dArray.Get(6)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Get 6 index", idx6)
	}
	

	fmt.Println("Length", dArray.Len())
	fmt.Println("Capacity", dArray.Cap())
	fmt.Print("Array list\n")
	dArray.Print()
}
