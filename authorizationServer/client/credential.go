package client

var Credentials = []Credential{
	{
		ClientId:     "client1",
		ClientSecret: "secret1",
		RedirectUrls: []string{"http://localhost:9001/callback"},
	},
}

type Credential struct {
	ClientId     string
	ClientSecret string
	RedirectUrls []string
}

func (c *Credential) RedirectUrlContains(url string) bool {
	for _, u := range c.RedirectUrls {
		if u == url {
			return true
		}
	}
	return false
}

func Find(clientId string) *Credential {
	for _, c := range Credentials {
		if c.ClientId == clientId {
			return &c
		}
	}
	return nil
}
