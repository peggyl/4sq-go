package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
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
	if code == "" {
		http.Error(w, "missing code, could not get token", 502)
	}

	token, err := client.GetAccessToken(code)
	if err != nil {
		http.Error(w, err.Error(), 502)
	}

	// HACK: quick api call to make sure we got the right token
	// TODO: refactor once the user resource is implemented.
	// TODO: support date versioning and mode
	u := fmt.Sprintf("https://api.foursquare.com/v2/users/self?oauth_token=%s&v=20160124&m=foursquare", token)
	resp, _ := http.Get(u)

	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	type user struct {
		FirstName string `json:"firstName"`
		LastName  string `json:"lastName"`
	}
	type result struct {
		Response struct {
			User user `json:"user"`
		} `json:"response"`
	}
	var res result
	json.Unmarshal(body, &res)
	fmt.Fprintf(w, "Hello, %s %s!\n", res.Response.User.FirstName, res.Response.User.LastName)
}

func main() {
	http.HandleFunc("/", mainHandler)
	http.HandleFunc("/redirect", redirectHandler)
	http.ListenAndServe(":8080", nil)
}
