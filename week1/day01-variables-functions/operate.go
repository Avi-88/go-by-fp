package main

import (
	"math"
	"errors"
)

type Operator string 

func Calculate (A int , B int, operator Operator) (float64, error) {
	switch operator {
	case "+":
		return add(A,B)
	case "-":
		return subtract(A,B)
	case "*":
		return multiply(A,B)
	case "/":
		return divide(A,B)
	case "%":
		return mod(A,B)
	case "^":
		return power(A,B)
	}

	return 0, errors.New("Invalid operation: the mentioned operation is not supported")
}

func add (A int, B int) (float64, error) {
	return float64(A+B), nil
}

func subtract (A int, B int) (float64, error) {
	return float64(A-B), nil
}

func multiply (A int, B int) (float64, error) {
	return float64(A*B), nil
}

func divide (A int, B int) (float64, error) {
	if B == 0 {
		return 0, errors.New("Invalid operation: cannot divide by 0")
	}
	return float64(A)/float64(B), nil
}

func mod (A int, B int) (float64, error) {
	if B == 0 {
		return 0, errors.New("Invalid operation: cannot mod by 0")
	}
	return float64(A%B), nil
}

func power (A int, B int) (float64, error) {
	return math.Pow(float64(A),float64(B)), nil
}