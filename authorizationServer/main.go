package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, World!")
	})

	fmt.Println("Server is running on port 8888...")
	if err := http.ListenAndServe(":8888", nil); err != nil {
		fmt.Printf("Error starting server: %s\n", err)
	}
}
