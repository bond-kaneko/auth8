package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/", Index)

	fmt.Println("Server is running on http://localhost:9001")
	http.ListenAndServe(":9001", nil)
}
