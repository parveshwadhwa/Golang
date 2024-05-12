package main

import "fmt"

func main() {
	fmt.Println("Welcome To Pointers")

	// var ptr *int

	// fmt.Println("value of pointer is: ", ptr)

	myNumber := 23

	var ptr = &myNumber
	fmt.Println("value of actual pointer is: ", ptr)  // reference to direct memory location
	fmt.Println("value of actual pointer is: ", *ptr) // value at that address

	*ptr = *ptr * 4

	fmt.Println("New value is: ", myNumber)
}
