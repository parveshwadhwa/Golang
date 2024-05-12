package main

import "fmt"

const LoginToken string = "ghksbdfud72fd" // if a variable first letter is capital , then it will be a public variable and accessile anywhere

func main() {
	var username string = "Prince"
	fmt.Println(username)
	fmt.Printf("Variable is of type is: %T \n", username)

	var isLogged bool = true
	fmt.Println(isLogged)
	fmt.Printf("Variable is of type is: %T \n", isLogged)

	var smallVal int = 256
	fmt.Println(smallVal)
	fmt.Printf("Variable is of type is: %T \n", smallVal)

	var smallValFloat float32 = 256.85544325
	fmt.Println(smallValFloat)
	fmt.Printf("Variable is of type is: %T \n", smallValFloat) // it will give uptoo 5th precision decimal point

	var smallValFloat2 float64 = 256.85544325
	fmt.Println(smallValFloat2)
	fmt.Printf("Variable is of type is: %T \n", smallValFloat2) // it is more precise

	// Deafult Values and Some aliases

	var anotherVariable int
	fmt.Println(anotherVariable)
	fmt.Printf("Variable is of type is: %T \n", anotherVariable)

	// implicit declaration

	var website = "code.com"
	fmt.Println(website)
	fmt.Printf("Variable is of type is: %T \n", website)

	// website = 3  this cant be allowed because if we didnt give any dataType then lexer will considered it upon the value given and here it is string , so it cant be changd to int further

	// No var style

	numberOfUsers := 300000
	fmt.Println(numberOfUsers)
	fmt.Printf("Variable is of type is: %T \n", numberOfUsers)

	fmt.Println(LoginToken)
	fmt.Printf("Variable is of type is: %T \n", LoginToken)

}

// jwtToken := 400 global deceleration is not possible using this := operator
