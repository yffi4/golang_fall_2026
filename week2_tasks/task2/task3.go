package main

import "fmt"

func modifySlice(s []int) {
	for i := range s {
		s[i] *= 10
	}
	fmt.Println(s)
}

func main() {
	numbers := []int{1, 2, 3, 4, 5}

	fmt.Println("Before calling function:", numbers)

	modifySlice(numbers)

	fmt.Println("After calling function:  ", numbers)
}
