package main

import (
	"encoding/json"
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
		e := struct {
			Error string
		}{
			Error: err.Error(),
		}
		jerr, err := json.Marshal(e)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Write(jerr)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	res := struct {
		Message string `json:"message"`
	}{
		Message: "Protected Resource",
	}
	resBody, err := json.Marshal(res)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, err)
		return
	}
	w.Write(resBody)
}
