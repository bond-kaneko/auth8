package main

import (
	"fmt"
	"net/http"

	"github.com/bond-kaneko/auth8/protectedResource/auth"
)

func Public(w http.ResponseWriter, r *http.Request) {
	html := `
		<!DOCTYPE html>
		<html>
		<head>
			<title>Protected Resource</title>
		</head>
		<body>
			<h1>Public Resource</h1>
		</body>
		</html>
	`
	fmt.Fprint(w, html)
}

func Protected(w http.ResponseWriter, r *http.Request) {
	_, err := auth.GetAccessToken(r)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprint(w, err)
		return
	}

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
	fmt.Fprint(w, html)
}
