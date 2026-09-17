package main

import "fmt"

func main() {
	// 1. Using a slice literal (like an array literal, but no length)
	s1 := []int{1, 2, 3}

	// 2. Using make() — specify length (and optionally capacity)
	s2 := make([]int, 5) // length 5, capacity 5, all zero values

	// 3. Slicing an existing array
	arr := [5]int{10, 20, 30, 40, 50}
	s3 := arr[1:4] // elements at index 1,2,3 -> [20 30 40]

	// 4. Declaring a nil slice (zero value of a slice type), then appending
	var s4 []int
	s4 = append(s4, 100, 200)

	fmt.Println("s1:", s1)
	fmt.Println("s2:", s2)
	fmt.Println("s3:", s3)
	fmt.Println("s4:", s4)
}
