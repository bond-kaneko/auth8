package auth

var Client = OAuthClient{
	Id:          "client1",
	Secret:      "secret1",
	RedirectUrl: "http://localhost:9001/callback",
}

type OAuthClient struct {
	Id          string
	Secret      string
	RedirectUrl string
}

var Server = OAuthServer{
	AuthorizationEndpoint: "http://localhost:9000/authorize",
	TokenEndpoint:         "http://localhost:9000/token",
}

type OAuthServer struct {
	AuthorizationEndpoint string
	TokenEndpoint         string
}
