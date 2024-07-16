package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// HTMLコンテンツを定義
		html := `
            <!DOCTYPE html>
            <html>
            <head>
                <title>Protected Resource</title>
            </head>
            <body>
                <h1>Protected Resource</h1>
            </body>
            </html>
        `
		// レスポンスとしてHTMLを返す
		fmt.Fprint(w, html)
	})

	// 8080ポートでサーバーを起動
	fmt.Println("Server is running on http://localhost:9002")
	http.ListenAndServe(":9002", nil)
}
