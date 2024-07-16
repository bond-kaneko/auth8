package client

var Credentials = []Credential{
	{
		ClientId:     "client1",
		ClientSecret: "secret1",
		CallbackUrls: []string{"http://localhost:9000/callback"},
	},
}

type Credential struct {
	ClientId     string
	ClientSecret string
	CallbackUrls []string
}
