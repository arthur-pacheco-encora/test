package main

import (
	"fmt"
	"net/http"
)

func helloWorldHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello, Encora!\n")
}

func main() {
	http.HandleFunc("/", helloWorldHandler)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(fmt.Sprintf("Server failed to start: %v\n", err))
	}
}
