package main

import "fmt"

func main() {
	fmt.Println("Maps in Golang")

	languages := make(map[string]string) // [string] for data type of key

	languages["JS"] = "javascript"
	languages["RB"] = "ruby"
	languages["PY"] = "python"

	fmt.Println("list of all languages: ", languages)
	fmt.Println("JS shorts for: ", languages["JS"])

	delete(languages, "RB")
	fmt.Println("List of languages:", languages)

	//  loops are interesting in golang

	for key, value := range languages {
		fmt.Printf("For Key %v, value is %v\n", key, value)
	}

	// for _, value := range languages { // for comma,ok syntax
	// 	fmt.Printf("For Key v, value is %v\n", value)
	// }

}
