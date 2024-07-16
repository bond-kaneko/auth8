package auth

var Client = OAuthClient{
	Id:          "client1",
	Secret:      "secret1",
	CallbackUrl: "http://localhost:9000/callback",
}

type OAuthClient struct {
	Id          string
	Secret      string
	CallbackUrl string
}

var Server = OAuthServer{
	AuthorizationEndpoint: "http://localhost:8888/authorize",
	TokenEndpoint:         "http://localhost:8888/token",
}

type OAuthServer struct {
	AuthorizationEndpoint string
	TokenEndpoint         string
}
