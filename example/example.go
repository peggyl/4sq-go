package main

import (
	"fmt"
	"net/http"
	"os"
	"text/template"

	foursquare "github.com/peggyl/4sq-go"
)

var (
	client foursquare.FoursquareClient
)

func init() {
	// First, go to developers.foursquare.com and create a new app.
	// This will give you a client ID and secret to use for OAuth.
	// Make sure your redirect URI matches one of the redirect URIs
	// that you specified when configuring your app with Foursquare.

	clientID := os.Getenv("FOURSQUARE_CLIENT_ID")
	clientSecret := os.Getenv("FOURSQUARE_CLIENT_SECRET")
	redirectURI := os.Getenv("FOURSQUARE_REDIRECT_URI")

	client = foursquare.NewClient(clientID, clientSecret, redirectURI)
	fmt.Printf("+v\n", client)
}

// mainHandler is the handler for the root URL (index).
// It will return a simple page containing a link to the
// authorization URL where the user authorizes the app.
func mainHandler(w http.ResponseWriter, r *http.Request) {
	authURL := client.AuthURL()
	tmpl, err := template.New("index").Parse("<!doctype html><html><body><a href=\"{{ .AuthURL }}\">Authorize Foursquare</a></body></html>")
	if err != nil {
		http.Error(w, err.Error(), 500)
	}
	data := struct{ AuthURL string }{authURL}
	tmpl.Execute(w, data)
}

// redirectHandler handles the redirect once an access token
// has been acquired from the Foursquare OAuth flow.
// In this example, we will return a page listing the name
// and current city (if available) of the now-logged-in user.
func redirectHandler(w http.ResponseWriter, r *http.Request) {
	// check for `code` param in URL
	code := r.URL.Query().Get("code")

	if code != "" {
		token, err := client.GetAccessToken(code)
		if err != nil {
			http.Error(w, err.Error(), 502)
		}

		fmt.Println(token)
	}

	// otherwise, parse the access token, woot.
}

func main() {
	http.HandleFunc("/", mainHandler)
	http.ListenAndServe(":8080", nil)
}
