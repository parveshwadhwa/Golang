package main

import "fmt"

func main() {
	fmt.Println("Structs in Golang")

	// There is no inheritance  in golang : No super or parent or child or anthing else like that

	prince := User{"Prince", "prince.com", true, 19}
	fmt.Println(prince)
	fmt.Printf("Prince Details are: %+v\n", prince)
	fmt.Printf("Name is %v and email is %v\n", prince.Name, prince.Email)

}

type User struct { // First letter of User is capital or  Name or Email or Status or Age , just to export them out and can be used anywhere
	Name   string
	Email  string
	Status bool
	Age    int
}
