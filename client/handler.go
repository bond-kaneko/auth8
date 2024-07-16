package main

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"
	"net/url"

	"github.com/bond-kaneko/auth8/client/auth"
	"github.com/bond-kaneko/auth8/client/internal"
)

func Index(w http.ResponseWriter, r *http.Request) {
	err := internal.Tmpl.ExecuteTemplate(w, "index.html", struct {
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
}

func GetAccessToken(w http.ResponseWriter, r *http.Request) {
	AccessToken = ""
	RefreshToken = ""

	url, err := authorizationEndpointUrl()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, url, http.StatusFound)
}

func authorizationEndpointUrl() (string, error) {
	authorizationUrl, err := url.Parse(auth.Server.AuthorizationEndpoint)
	if err != nil {
		return "", err
	}

	q := authorizationUrl.Query()
	q.Set("response_type", "code")
	q.Set("client_id", auth.Client.Id)
	q.Set("client_secret", auth.Client.Secret)
	q.Set("redirect_uri", auth.Client.RedirectUrl)
	q.Set("state", mustRandomString(10))
	authorizationUrl.RawQuery = q.Encode()

	return authorizationUrl.String(), nil
}

func mustRandomString(n int) string {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}
	return base64.URLEncoding.EncodeToString(b)
}
