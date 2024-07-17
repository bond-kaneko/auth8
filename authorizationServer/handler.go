package main

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/url"
	"strings"

	"github.com/bond-kaneko/auth8/client"
	"github.com/bond-kaneko/auth8/internal"
)

var AuthorizeRequests = map[string]AuthorizeRequest{}

type AuthorizeRequest struct {
	ClientId     string
	ResponseType string
	RedirectUrl  string
	State        string
}

func Authorize(w http.ResponseWriter, r *http.Request) {
	clientId := r.URL.Query().Get("client_id")
	if clientId == "" {
		http.Error(w, "client_id is required", http.StatusBadRequest)
		return
	}
	foundClient := client.Find(clientId)
	if foundClient == nil {
		http.Error(w, "client_id is invalid", http.StatusBadRequest)
		return
	}

	redirectUrl := r.URL.Query().Get("redirect_uri")
	if redirectUrl == "" {
		http.Error(w, "redirect_uri is required", http.StatusBadRequest)
		return
	}
	if !foundClient.RedirectUrlContains(redirectUrl) {
		http.Error(w, "redirect_uri is invalid", http.StatusBadRequest)
		return
	}

	state := r.URL.Query().Get("state")
	if state == "" {
		http.Error(w, "state is required", http.StatusBadRequest)
		return
	}

	// インメモリストアへ保存する
	requestId := mustRandomString(10)
	AuthorizeRequests[requestId] = AuthorizeRequest{
		ClientId:     clientId,
		RedirectUrl:  redirectUrl,
		ResponseType: "code",
		State:        state,
	}

	err := internal.Tmpl.ExecuteTemplate(w, "login.html", struct {
		RequestId string
		UserName  string
		UserId    string
	}{
		RequestId: requestId,
		UserName:  "user1",
		UserId:    "1",
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func mustRandomString(n int) string {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}
	return strings.ReplaceAll(base64.URLEncoding.EncodeToString(b), "=", "")
}

type CodeRequest struct {
	Code             string
	AuthorizeRequest AuthorizeRequest
	User             User
}

var CodeRequests = map[string]CodeRequest{}

func Approve(w http.ResponseWriter, r *http.Request) {
	userId := r.FormValue("userId")
	if userId == "" {
		http.Error(w, "userId is required", http.StatusBadRequest)
		return
	}
	userName := r.FormValue("userName")
	if userName == "" {
		http.Error(w, "userName is required", http.StatusBadRequest)
		return
	}
	requestId := r.FormValue("requestId")
	if requestId == "" {
		http.Error(w, "requestId is required", http.StatusBadRequest)
		return
	}
	authReq, ok := AuthorizeRequests[requestId]
	if !ok {
		http.Error(w, "requestId is invalid", http.StatusBadRequest)
		return
	}
	delete(AuthorizeRequests, requestId)

	code := mustRandomString(10)
	CodeRequests[code] = CodeRequest{
		Code:             code,
		AuthorizeRequest: authReq,
		User: User{
			Id:   userId,
			Name: userName,
		},
	}

	redirectUrl, err := url.Parse(authReq.RedirectUrl)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	q := redirectUrl.Query()
	q.Set("code", code)
	q.Set("state", authReq.State)
	redirectUrl.RawQuery = q.Encode()

	http.Redirect(w, r, redirectUrl.String(), http.StatusFound)
}

var IssuedTokens = []Token{}

type Token struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	User        User   `json:"user"`
}

func IssueToken(w http.ResponseWriter, r *http.Request) {
	auth := r.Header.Get("Authorization")
	if auth == "" {
		http.Error(w, "Authorization header is required", http.StatusBadRequest)
		return
	}
	credentialsBytes, err := base64.StdEncoding.DecodeString(auth[len("basic "):])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	credentials := string(credentialsBytes)
	clientId := strings.Split(credentials, ":")[0]

	clientSecret := strings.Split(credentials, ":")[1]
	cl := client.Find(clientId)
	if cl == nil {
		http.Error(w, "client_id is invalid", http.StatusBadRequest)
		return
	}
	if cl.ClientSecret != clientSecret {
		http.Error(w, "client_secret is invalid", http.StatusBadRequest)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	grantType := r.PostForm.Get("grant_type")
	if grantType != "authorization_code" {
		http.Error(w, "grant_type is invalid", http.StatusBadRequest)
		return
	}

	code := r.Form.Get("code")
	codeReq, ok := CodeRequests[code]
	if !ok {
		http.Error(w, "code is invalid", http.StatusBadRequest)
		return
	}
	delete(CodeRequests, code)

	if codeReq.AuthorizeRequest.ClientId != clientId {
		http.Error(w, "client_id is invalid", http.StatusBadRequest)
		return
	}

	accessToken := mustRandomString(10)
	IssuedTokens = append(IssuedTokens, Token{
		AccessToken: accessToken,
		TokenType:   "bearer",
		User:        codeReq.User,
	})

	w.Header().Set("Content-Type", "application/json")
	resBody := struct {
		AccessToken string `json:"access_token"`
		TokenType   string `json:"token_type"`
	}{
		AccessToken: accessToken,
		TokenType:   "bearer",
	}
	resJson, err := json.Marshal(resBody)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(resJson)
}
