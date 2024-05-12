package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

func main() {
	fmt.Println("Welcome to Go Requests")
	// PerformGetRequest()
	// PerformPostJSONRequest()
	PerformPOstFormDataRequest()
}

func PerformGetRequest() {
	const myUrl = "http://localhost:3000/get"

	response, err := http.Get(myUrl)

	if err != nil {
		panic(err)
	}

	defer response.Body.Close()

	fmt.Println("Status Code: ", response.StatusCode)
	fmt.Println("Content Length is: ", response.ContentLength)

	var responseString strings.Builder // string is a datatype but strings is a package

	content, err := ioutil.ReadAll(response.Body) // To read content of body and content is in byte format
	byteCount, _ := responseString.Write(content)

	fmt.Println("Byte count is: ", byteCount)
	// fmt.Println(string(content))
	fmt.Println(responseString.String())
	// A builder is used to efficiently build a string using Write methods. it minimizes memory copying. The zero value is ready to use. Do not copy a zero builder
}

func PerformPostJSONRequest() {
	const myUrl = "http://localhost:3000/post"

	requestBody := strings.NewReader(`
	     {
			"coursensme":"Lets go with golang",
			"price":"0",
			"platform":"leadcode.in"
		 }
	`)

	response, err := http.Post(myUrl, "application/json", requestBody)

	if err != nil {
		panic(err)
	}

	defer response.Body.Close()
	content, _ := ioutil.ReadAll(response.Body)

	var readString strings.Builder
	readString.Write(content)

	fmt.Println(readString.String())
}

func PerformPOstFormDataRequest() {
	const myUrl = "http://localhost:3000/postform"

	//formdata

	data := url.Values{}
	data.Add("firstName", "Parvesh")
	data.Add("LastName", "Wadhwa")
	data.Add("Email", "prince@gmail.com")

	response, err := http.PostForm(myUrl, data)

	if err != nil {
		panic(err)
	}
	defer response.Body.Close()
	content, _ := ioutil.ReadAll((response.Body))

	fmt.Println(string(content))

}
