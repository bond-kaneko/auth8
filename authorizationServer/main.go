package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/authorize", Authorize)
	http.HandleFunc("/approve", Approve)

	fmt.Println("Server is running on port 9000...")
	if err := http.ListenAndServe(":9000", nil); err != nil {
		fmt.Printf("Error starting server: %s\n", err)
	}
}
