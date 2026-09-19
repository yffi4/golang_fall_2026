package main

import "fmt"

func main() {

	s1 := []int{1, 2, 3}


	s2 := make([]int, 5)

	arr := [5]int{10, 20, 30, 40, 50}
	s3 := arr[1:4]


	var s4 []int
	s4 = append(s4, 100, 200)

	fmt.Println("s1:", s1)
	fmt.Println("s2:", s2)
	fmt.Println("s3:", s3)
	fmt.Println("s4:", s4)
}
