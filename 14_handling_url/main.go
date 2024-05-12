package main

import (
	"fmt"
	"net/url"
)

const myUrl string = "https://lco.dev:3000/learn?coursename=reactjs&paymentid=pr90"

func main() {
	fmt.Println("Welcome to handling url's in golang")
	fmt.Println(myUrl)

	// parsing
	result, err := url.Parse(myUrl)

	if err != nil {
		panic(err)
	}

	fmt.Println(result.Scheme)   // https
	fmt.Println(result.Host)     // lco.dev:3000
	fmt.Println(result.Path)     // /learn
	fmt.Println(result.Port())   // 3000
	fmt.Println(result.RawQuery) // coursename=reactjs&paymentid=pr90 (i.e params or parameters)

	qParams := result.Query()

	fmt.Printf("The type of query Params is %T\n", qParams) // key value pairs

	fmt.Println(qParams) // map[coursename:[reactjs] paymentid:[pr90]]

	fmt.Println(qParams["coursename"])

	for _, val := range qParams {
		fmt.Println("Param is : ", val)
	}

	partsOfUrl := &url.URL{
		Scheme:   "https",
		Host:     "prince.dev",
		Path:     "/go",
		RawQuery: "user=parvesh",
	}

	anotherUrl := partsOfUrl.String()

	fmt.Println(anotherUrl)

}
