# Day 3 — Structs, Methods, and Interfaces

## Goal
Understand how Go models behavior through interfaces — not inheritance.
This is how Go replaces class hierarchies. It's a fundamentally different mental model.

---

## The Exercise: Shape Calculator

Build a shape area/perimeter calculator using interfaces.

---

## Part 1: Define the Interface

```go
type Shape interface {
    Area() float64
    Perimeter() float64
    String() string   // human-readable name + dimensions
}
```

---

## Part 2: Implement 3 Concrete Shapes

### Circle
```go
type Circle struct {
    Radius float64
}
```

### Rectangle
```go
type Rectangle struct {
    Width  float64
    Height float64
}
```

### Triangle (right triangle)
```go
type Triangle struct {
    Base   float64
    Height float64
}
```

Implement `Area()`, `Perimeter()`, and `String()` for all three.

Formulas:
- Circle area: `π * r²`, perimeter (circumference): `2 * π * r`
- Rectangle area: `w * h`, perimeter: `2 * (w + h)`
- Right triangle area: `0.5 * b * h`, perimeter: `b + h + √(b² + h²)`

Use `math.Pi` and `math.Sqrt` from the `math` package.

---

## Part 3: A Function That Accepts Any Shape

Write this function:

```go
func PrintShapeInfo(s Shape) {
    // print the shape's String(), Area(), and Perimeter()
}
```

Call it from `main` with all three shapes. Notice: `PrintShapeInfo` doesn't know or care what concrete type it receives — it only knows it has `Area()`, `Perimeter()`, and `String()`.

---

## Part 4: A Slice of Shapes

```go
shapes := []Shape{
    Circle{Radius: 5},
    Rectangle{Width: 4, Height: 6},
    Triangle{Base: 3, Height: 4},
}
```

Loop over the slice and call `PrintShapeInfo` on each. Then find and print the shape with the largest area.

---

## Part 5: Type Assertion

Write a function:

```go
func Describe(s Shape) {
    // use a type switch to print something specific
    // about the concrete type, not just the interface
}
```

Use a **type switch**:

```go
switch v := s.(type) {
case Circle:
    fmt.Printf("This is a circle with radius %.2f\n", v.Radius)
case Rectangle:
    // ...
case Triangle:
    // ...
}
```

Call `Describe` on each shape from `main`.

---

## Concepts to understand deeply

- A type satisfies an interface **implicitly** — no `implements` keyword
- If a type has all the methods the interface requires, it satisfies it automatically
- The interface variable holds two things internally: `{type, value}` — called a "fat pointer"
- A `nil` interface is different from an interface holding a `nil` pointer (this causes bugs)

---

## Intentional breaks (do these after it works)

1. Remove `Area()` from one of your shapes. Read the compiler error — Go tells you exactly which method is missing.
2. Try assigning a `Circle` directly to a `Shape` variable: `var s Shape = Circle{Radius: 5}`. It works. Now try: `var s Shape = &Circle{Radius: 5}`. Does it still work? Why?
3. Add a fourth method `Color() string` to your `Shape` interface without implementing it on any shape. Watch the compiler errors cascade.
4. Try a type assertion without a type switch: `c := s.(Circle)`. Then try `c := s.(Rectangle)` when `s` is actually a `Circle`. See the panic. Then use the safe form: `c, ok := s.(Rectangle)`.

---

## When done, write NOTES.md

- How is a Go interface different from an abstract class in another language you know?
- What does "implicit satisfaction" mean and why is it useful?
- What are the two things stored inside an interface variable?
- What's the difference between `s.(Circle)` and `s.(type)` in a type switch?
