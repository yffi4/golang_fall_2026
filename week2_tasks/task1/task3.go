package main

import "fmt"

func main() {
	var a [3][3]int
	a[0][0] = 1
	a[1][1] = 2
	a[2][2] = 3

	fmt.Println(a)
}
