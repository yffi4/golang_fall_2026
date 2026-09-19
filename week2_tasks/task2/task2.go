package main

import "fmt"

func main() {
	arr := [5]int{1, 2, 3, 4, 5}
	s := arr[1:4]

	fmt.Println("Before:")
	fmt.Println("arr:", arr)
	fmt.Println("s:  ", s)


	s[0] = 99
	fmt.Println("\nAfter modifying s[0] = 99:")
	fmt.Println("arr:", arr)
	fmt.Println("s:  ", s)


	arr[2] = 500
	fmt.Println("\nAfter modifying arr[2] = 500:")
	fmt.Println("arr:", arr)
	fmt.Println("s:  ", s)
}
