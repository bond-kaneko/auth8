package main

import (
	"fmt"
	"net/http"

	"github.com/bond-kaneko/auth8/client/auth"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		html := `
            <!DOCTYPE html>
            <html>
            <head>
                <title>OAuth Client</title>
            </head>
            <body>
				<h1>OAuth Client</h1>
                <div style="text-align: center; margin: auto; width: 50%%;">
					<div>
						<h2>Client</h2>
						<table style="text-align: left">
							<tr>
								<th>Client ID</th>
								<td>%s</td>
							</tr>
							<tr>
								<th>Client Secret</th>
								<td>%s</td>
							</tr>
							<tr>
								<th>Callback URL</th>
								<td>%s</td>
							</tr>
						</table>
					</div>
					<div>
						<h2>Authorization Server</h2>
						<table style="text-align: left">
							<tr>
								<th>Authorization Endpoint</th>
								<td>%s</td>
							</tr>
							<tr>
								<th>Token Endpoint</th>
								<td>%s</td>
							</tr>
						</table>
					</div>
				</div>
            </body>
            </html>
        `
		fmt.Fprintf(w, html, auth.Client.Id, auth.Client.Secret, auth.Client.CallbackUrl, auth.Server.AuthorizationEndpoint, auth.Server.TokenEndpoint)
	})

	fmt.Println("Server is running on http://localhost:9000")
	http.ListenAndServe(":9000", nil)
}
