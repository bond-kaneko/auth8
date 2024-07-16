package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/", Index)
	http.HandleFunc("/accessToken", GetAccessToken)
	http.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "error: %s, code: %s", r.URL.Query().Get("error"), r.URL.Query().Get("code"))
	})

	fmt.Println("Server is running on http://localhost:9001")
	http.ListenAndServe(":9001", nil)
}
