package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

// net/http is the recommnded package for handling web requests and to learn about it we can learn from official docs

const url = "https://lco.dev"

func main() {
	fmt.Println("Welcome to LCO Web Request")

	response, err := http.Get(url)

	if err != nil {
		panic(err)
	}

	fmt.Printf("Response is of Type: %T\n", response) // output will be *http.Response , it means we are getting the origional response not the copy of response as we are getting a pointers and we can further manipulate this reponse
	defer response.Body.Close()                       // caller's responsibilty to close the connection

	dataBytes, err := ioutil.ReadAll(response.Body)

	if err != nil {
		panic(err)
	}

	content := string(dataBytes)

	fmt.Println(content)
}
