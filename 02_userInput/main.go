package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	welcome := "Welcome to User Input"
	fmt.Println(welcome)

	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Enter:")

	// comma ok / comma err syntax

	input, _ := reader.ReadString('\n')
	fmt.Println("Thanks for reading ", input)
	fmt.Printf("Type of reading is %T ", input)

	// we can use _, err also input, err

}
