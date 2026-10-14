package main

import "fmt"

type Car struct {
	brand  string
	ID     int
	model  string
	engine float64
	colors []string
}

func addColor(c Car, color string) {
	c.colors = append(c.colors, color)

}
func addColorPtr(c *Car, color string) {
	c.colors = append(c.colors, color)
}
func main() {

	var c1 Car
	c2 := Car{brand: "Toyota", ID: 1, model: "Land Cruiser Prado 120", engine: 2.7}
	c3 := Car{"Mersedes", 2, "E211", 3.2, []string{"blue", "red"}}
	p := &Car{brand: "BMW"}
	addColor(c2, "blue")
	addColorPtr(&c1, "blue")

	fmt.Printf("%+v\n", c1)
	fmt.Printf("%+v\n", c2)
	fmt.Printf("%+v\n", c3)
	fmt.Printf("%+v\n", p)
}
