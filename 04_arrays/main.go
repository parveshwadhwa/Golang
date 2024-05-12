package main

import "fmt"

/*
We mostly use slices instead of arrays
*/

func main() {
	fmt.Println("Welcome to array")

	var fruitList [4]string

	fruitList[0] = "Apple"
	fruitList[1] = "Mango"
	fruitList[3] = "peach"

	fmt.Println("Fruitlist is: ", fruitList)
	fmt.Println("Fruitlist is: ", len(fruitList))

	var vegetableList = [5]string{"potato", "beans", "mushroom"}
	fmt.Println("Veggie List is: ", vegetableList)
	fmt.Println("Veggie List is: ", len(vegetableList))

}
