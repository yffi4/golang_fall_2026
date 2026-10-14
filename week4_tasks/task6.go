package main

import (
	"fmt"
	"math"
)

type Shape interface {
	Area() float64
	Perimeter() float64
}

type Rectangle struct {
	Width, Height float64
}

func (r Rectangle) Area() float64 {
	return r.Height * r.Width
}
func (r Rectangle) Perimeter() float64 {
	return 2 * (r.Width + r.Height)
}

func (r *Rectangle) Scale(k float64) {
	r.Width *= k
	r.Height *= k
}

type Circle struct {
	Radius float64
}

func (c *Circle) Area() float64 {
	return math.Pi * c.Radius * c.Radius
}

func (c *Circle) Perimeter() float64 {
	return 2 * math.Pi * c.Radius
}

func main() {
	rect := Rectangle{Width: 4, Height: 5}
	rect.Scale(2)
	fmt.Printf("Rectangle after scale: %+v\n", rect)

	circ := &Circle{Radius: 3}


	shapes := []Shape{rect, circ}

	for _, s := range shapes {
		fmt.Printf("Area: %.2f, Perimeter: %.2f\n", s.Area(), s.Perimeter())
	}

}
