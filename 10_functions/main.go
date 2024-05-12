package main

import "fmt"

func main() {
	fmt.Println("Welcome to functions")
	greeter()

	result := adder(3, 5)
	result2, myMessage := proAdder(3, 5, 6, 7, 89)

	fmt.Println("Result is: ", result)
	fmt.Println("Result2 is: ", result2)
	fmt.Println("Message is: ", myMessage)
}

func adder(valOne int, valTwo int) int {
	return valOne + valTwo
}

func proAdder(values ...int) (int, string) {
	total := 0

	for _, value := range values {
		total += value
	}

	return total, "Hii prince"
}

func greeter() {
	fmt.Println("Hii Prince namastey")
}
