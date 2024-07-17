package main

import (
	"fmt"
	"net/http"

	"github.com/rs/cors"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/authorize", Authorize)
	mux.HandleFunc("/approve", Approve)
	mux.HandleFunc("/token", IssueToken)

	fmt.Println("Server is running on port 9000...")
	if err := http.ListenAndServe(":9000", cors.Default().Handler(mux)); err != nil {
		fmt.Printf("Error starting server: %s\n", err)
	}
}
