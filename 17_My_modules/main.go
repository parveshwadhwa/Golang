package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	fmt.Println("Hello Mod In Golang")
	greeter()

	mux := mux.NewRouter()
	mux.HandleFunc("/", serverHome).Methods("GET")
	log.Fatal(http.ListenAndServe(":4000", mux))
}

func greeter() {
	fmt.Println("Hello Mod in Golang")
}

func serverHome(w http.ResponseWriter, r *http.Request) { // The URL or the params we get or want will come under r *http.Request and if we want to some send response then it will done through w http.ResponseWriter
	w.Write([]byte("<h1>Welcome To Golang Series</h1>"))
}
