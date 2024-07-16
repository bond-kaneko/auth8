package main

import (
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
)

func mustLoadTemplates(pattern string) *template.Template {
	// filepath.Globを使用してマッチする全てのファイルのパスを取得
	files, err := filepath.Glob(pattern)
	if err != nil {
		panic(err)
	}

	// template.ParseFilesにファイルのリストを渡してテンプレートをパース
	tmpl, err := template.ParseFiles(files...)
	if err != nil {
		panic(err)
	}
	return tmpl
}

var tmpl = mustLoadTemplates("web/html/*")

func main() {

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		err := tmpl.ExecuteTemplate(w, "index.html", struct {
			ClientId              string
			ClientSecret          string
			CallbackUrl           string
			AuthorizationEndpoint string
			TokenEndpoint         string
		}{
			ClientId:              "client1",
			ClientSecret:          "secret1",
			CallbackUrl:           "http://localhost:9001/callback",
			AuthorizationEndpoint: "http://localhost:9000/authorize",
			TokenEndpoint:         "http://localhost:9000/token",
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	fmt.Println("Server is running on http://localhost:9001")
	http.ListenAndServe(":9001", nil)
}
