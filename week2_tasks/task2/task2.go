package main

import "fmt"

func main() {
	arr := [5]int{1, 2, 3, 4, 5}
	s := arr[1:4] // slice view into arr: [2 3 4]

	fmt.Println("Before:")
	fmt.Println("arr:", arr)
	fmt.Println("s:  ", s)

	// Modify through the slice
	s[0] = 99
	fmt.Println("\nAfter modifying s[0] = 99:")
	fmt.Println("arr:", arr) // arr also changes!
	fmt.Println("s:  ", s)

	// Modify through the array
	arr[2] = 500
	fmt.Println("\nAfter modifying arr[2] = 500:")
	fmt.Println("arr:", arr)
	fmt.Println("s:  ", s) // s also changes!
}
