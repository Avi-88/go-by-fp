package main

import (
	"math"
	"strconv"
	"fmt"
)
// Day 3: Structs, Methods, and Interfaces
//
// Read EXERCISE.md before starting.
// Implement all 5 parts in this file.

type Shape interface {
	Area() float64
	Perimeter() float64
	String() string
}

type Circle struct {
	Radius float64
}

type Rectangle struct {
	Width float64
	Height float64
}

type Triangle struct {
	Base float64
	Height float64
}

func (c Circle) Area () float64 {
	return math.Pi*(c.Radius*c.Radius)
}

func (r Rectangle) Area () float64 {
	return r.Width*r.Height
}

func (t Triangle) Area () float64 {
	return 0.5*t.Base*t.Height
}

func (c Circle) Perimeter () float64 {
	return 2*math.Pi*(c.Radius)
}

func (r Rectangle) Perimeter () float64 {
	return 2*(r.Width  + r.Height)
}

func (t Triangle) Perimeter () float64 {
	return t.Base + t.Height + math.Sqrt(math.Pow(t.Base, 2) + math.Pow(t.Height, 2))
}

func (c Circle) String () string {
	r := c.Radius
	return "Circle with radius " + strconv.FormatFloat(r, 'f', 2, 64)
}

func (r Rectangle) String () string {
	w := r.Width
	h := r.Height
	return "Rectangle " + strconv.FormatFloat(w,'f', 2, 64) + " x " + strconv.FormatFloat(h, 'f', 2, 64)
}

func (t Triangle) String () string {
	b := t.Base
	h := t.Height
	return "Triangle " + strconv.FormatFloat(b,'f', 2, 64) + " x " + strconv.FormatFloat(h, 'f', 2, 64)
}

func PrintShapeInfo(s Shape) {
    // print the shape's String(), Area(), and Perimeter()
	fmt.Printf("String - %s\n", s.String())
	fmt.Printf("Area - %f\n", s.Area())
	fmt.Printf("Perimeter - %f\n", s.Perimeter())
}

func Describe(s Shape) {
    // use a type switch to print something specific
    // about the concrete type, not just the interface
	switch v:= s.(type) {
	case Circle:
		fmt.Printf("This is a circle with radius %v\n", v.Radius)
		return 
	case Rectangle:
		fmt.Printf("This is a rectangle with height - %v and width - %v\n", v.Height, v.Width)
		return
	case Triangle:
		fmt.Printf("This is a triangle with base - %v and height - %v\n", v.Base, v.Height)
		return
	default:
		fmt.Printf("An unknown type was passed as a shape")
	}	
}

func main() {
	shapes := []Shape{
		Circle{Radius: 5},
		Rectangle{Width: 4, Height: 6},
		Triangle{Base: 3, Height: 4},
	}

	maxArea := -math.MaxFloat64
	maxAreaShape := "No valid shape found"
	for _,i := range shapes {
		PrintShapeInfo(i)
		if maxArea < i.Area() {
			maxArea = i.Area()
			maxAreaShape = i.String()
		}
		Describe(i)
	}
	fmt.Printf("The shape with the max area is: %s\n", maxAreaShape)
}
