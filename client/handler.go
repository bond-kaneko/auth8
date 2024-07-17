package main

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/bond-kaneko/auth8/client/auth"
	"github.com/bond-kaneko/auth8/client/internal"
)

func Index(w http.ResponseWriter, r *http.Request) {
	err := internal.Tmpl.ExecuteTemplate(w, "index.html", struct {
		AuthorizationCode     string
		AccessToken           string
		RefreshToken          string
		ClientId              string
		ClientSecret          string
		CallbackUrl           string
		AuthorizationEndpoint string
		TokenEndpoint         string
	}{
		AuthorizationCode:     AuthorizationCode,
		AccessToken:           AccessToken,
		RefreshToken:          RefreshToken,
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

func Authorize(w http.ResponseWriter, r *http.Request) {
	AccessToken = ""
	RefreshToken = ""
	AuthorizationCode = ""
	State = mustRandomString(10)

	url, err := authorizationEndpointUrl(State)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, url, http.StatusFound)
}

func authorizationEndpointUrl(state string) (string, error) {
	authorizationUrl, err := url.Parse(auth.Server.AuthorizationEndpoint)
	if err != nil {
		return "", err
	}

	q := authorizationUrl.Query()
	q.Set("response_type", "code")
	q.Set("client_id", auth.Client.Id)
	q.Set("client_secret", auth.Client.Secret)
	q.Set("redirect_uri", auth.Client.RedirectUrl)
	q.Set("state", state)
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

func Callback(w http.ResponseWriter, r *http.Request) {
	inState := r.URL.Query().Get("state")
	if inState != State {
		http.Error(w, "state is invalid", http.StatusBadRequest)
		return
	}

	AuthorizationCode = r.URL.Query().Get("code")
	values := url.Values{
		"grant_type":   {"authorization_code"},
		"code":         {AuthorizationCode},
		"redirect_uri": {auth.Client.RedirectUrl},
	}
	tokReq, err := http.NewRequest("POST", auth.Server.TokenEndpoint, strings.NewReader(values.Encode()))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tokReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	credential := base64.StdEncoding.EncodeToString([]byte(auth.Client.Id + ":" + auth.Client.Secret))
	tokReq.Header.Set("Authorization", "Basic "+credential)

	tokRes, err := http.DefaultClient.Do(tokReq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if tokRes.StatusCode != http.StatusOK {
		body, err := io.ReadAll(tokRes.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.Error(w, "failed to get token: "+string(body), http.StatusInternalServerError)
		return
	}

	defer tokRes.Body.Close()
	resBodyBytes, err := io.ReadAll(tokRes.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tokenRes := struct {
		AccessToken string `json:"access_token"`
		TokenType   string `json:"token_type"`
	}{}
	if err := json.Unmarshal(resBodyBytes, &tokenRes); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	AccessToken = tokenRes.AccessToken

	http.Redirect(w, r, "http://localhost:9001/index.html", http.StatusFound)
}

func GetProtectedResource(w http.ResponseWriter, r *http.Request) {
	if AccessToken == "" {
		http.Error(w, "access token is empty", http.StatusBadRequest)
		return
	}

	req, err := http.NewRequest("GET", "http://protected-resource:9002/protected", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	req.Header.Set("Authorization", "Bearer "+AccessToken)

	protectedRes, err := http.DefaultClient.Do(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if protectedRes.StatusCode != http.StatusOK {
		body, err := io.ReadAll(protectedRes.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.Error(w, "failed to get protected resource: "+string(body), http.StatusInternalServerError)
		return
	}

	var protectedResBody = struct {
		Message string `json:"message"`
	}{}
	defer protectedRes.Body.Close()
	err = json.NewDecoder(protectedRes.Body).Decode(&protectedResBody)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res, err := json.Marshal(protectedResBody)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = internal.Tmpl.ExecuteTemplate(w, "fetch_resource.html", struct {
		ProtectedResource string
	}{
		ProtectedResource: string(res),
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
