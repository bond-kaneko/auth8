package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/public", Public)

	http.HandleFunc("/protected", Protected)

	fmt.Println("Server is running on http://localhost:9002")
	http.ListenAndServe(":9002", nil)
}
