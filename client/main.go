package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/", Index)
	http.HandleFunc("/authorize", Authorize)
	http.HandleFunc("/callback", Callback)

	fmt.Println("Server is running on http://localhost:9001")
	http.ListenAndServe(":9001", nil)
}
