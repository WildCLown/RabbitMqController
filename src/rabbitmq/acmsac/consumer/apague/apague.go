package main

import "fmt"

var x int

func nothing(a *int) {
	*a++
}

func main() {
	x = 10
	nothing(&x)
	fmt.Println(x)
}
