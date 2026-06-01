package main

import (
	"fmt"
	"os"
	"strconv"
)

// Day 1: Variables, Types, Functions
//
// YOUR TASK: implement this calculator.
// Read EXERCISE.md first.
//
// Delete this comment block when you start writing.

func main() {
	userArgs := os.Args[1:]

	if len(userArgs) != 3 {
		fmt.Println("error: usage: calc <num> <op> <num>")
		return
	}
	
	var operation Operator = Operator(userArgs[1])
	num1, err := strconv.Atoi(userArgs[0])
	if err != nil {
		fmt.Println("Invalid input passed: num1 ")
		return
	}
	num2, err := strconv.Atoi(userArgs[2])
	if err != nil {
		fmt.Println("Invalid input passed")
		return
	}

	ans, err := Calculate(num1, num2, operation)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("%d %s %d = %g\n", num1, operation, num2, ans)
}
