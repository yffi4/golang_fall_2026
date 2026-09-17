package main

import "fmt"

func main() {
	a := [5]int{1, 2, 3}
	fmt.Println(a)
}

//It's because arrays in Go have a fixed length. When you use a composite literal like {1, 2, 3}
//and give fewer values than the declared length, Go doesn't leave the remaining
//slots undefined or throw an error — it fills them with the zero value for the element type.

