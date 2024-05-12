package main

import (
	"fmt"
	"sort"
)

func main() {
	fmt.Println("Welcome to slices")

	var fruitList = []string{"apple", "tomato", "peach"}
	fmt.Printf("Type of fruitlist is %T\n", fruitList)

	fruitList = append(fruitList, "mango", "banana")
	fmt.Println("FruilList is: ", fruitList)

	fruitList = append(fruitList[1:]) // this : syntax is used to slice your slice it means just to seperate slice into different parts , here it means apple is no longer availale , simply it means i will start from 0 and i will skip zero and ending limit is off
	fmt.Println(fruitList)

	// fruitList = append(fruitList[1:3]) // it means i will start from 0 and skip 0th index and ended upto 2nd leaving 3rd index
	// fmt.Println(fruitList)

	fruitList = append(fruitList[:3]) // it means starting from zero but ends at 2nd
	fmt.Println(fruitList)

	highScore := make([]int, 4) // heap memory

	highScore[0] = 234
	highScore[1] = 298
	highScore[2] = 236
	highScore[3] = 237
	// highScore[4] = 457

	highScore = append(highScore, 455, 459, 789)
	fmt.Println(sort.IntsAreSorted(highScore))
	sort.Ints(highScore)
	fmt.Println("Sorted highScores: ", highScore)
	fmt.Println(sort.IntsAreSorted(highScore))

	fmt.Println("highScores: ", highScore)

	// How To Remove a value from slices based on index

	var courses = []string{"react.js", "javascript", "python", "ruby", "go"}
	fmt.Println("Corses are: ", courses)

	var index int = 2
	// courses = append(courses[:index])
	// fmt.Println("List of courses: ", courses)

	// now if we want to remove value at index 2

	courses = append(courses[:index], courses[index+1:]...) //  it means first  list will be 0,1 and then second list will start from 3rd index
	fmt.Println("List of courses: ", courses)
}
