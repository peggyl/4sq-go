package foursquare

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

// FoursquareClient wraps authentication and API-querying logic with
// Foursquare.
type FoursquareClient struct {
	ClientID     string
	ClientSecret string
	AccessToken  string
	RedirectURI  string
}

// NewClient creates a new FoursquareClient for a registered application.
func NewClient(clientID string, clientSecret string, redirectURI string) FoursquareClient {
	//	if redirectURI == "" {
	//		return FoursquareClient{}, fmt.Errorf("must provide a redirect URI")
	//	}

	return FoursquareClient{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURI:  redirectURI,
	}
}

// AuthURL returns the URL where the user can authorize the app.
// See step 1: https://developer.foursquare.com/overview/auth
func (c FoursquareClient) AuthURL() string {
	u := url.URL{
		Scheme: "https",
		Path:   OAuthRoot + AuthEndpoint,
		RawQuery: url.Values{
			"client_id":     []string{c.ClientID},
			"response_type": []string{"code"},
			"redirect_uri":  []string{c.RedirectURI},
		}.Encode(),
	}
	return u.String()
}

// GetAccessToken returns the access code returned by Foursquare
// by exchanging the code received, or an error if one occurs.
// See step 3 of: https://developer.foursquare.com/overview/auth
func (c FoursquareClient) GetAccessToken(code string) (string, error) {
	if code == "" {
		return "", fmt.Errorf("must provide code")
	}
	u := url.URL{
		Scheme: "https",
		Path:   OAuthRoot + TokenEndpoint,
		RawQuery: url.Values{
			"client_id":     []string{c.ClientID},
			"client_secret": []string{c.ClientSecret},
			"grant_type":    []string{"authorization_code"},
			"redirect_uri":  []string{c.RedirectURI},
			"code":          []string{code},
		}.Encode(),
	}
	resp, err := http.Get(u.String())
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var token struct {
		AccessToken string `json:"access_token"`
	}
	if err = json.Unmarshal(body, &token); err != nil {
		return "", err
	}

	// set access token on the client and return it anyway (?)
	c.AccessToken = token.AccessToken
	return token.AccessToken, nil
}

func (c FoursquareClient) get(path string, params map[string]string) (*http.Response, error) {
	// all requests against the Foursquare API need a v(ersion), m(ode), and token
	vals := url.Values{}
	vals.Set("v", DateVersion)
	vals.Set("m", FoursquareMode)
	if c.AccessToken != "" {
		vals.Set("oauth_token", c.AccessToken)
	} else {
		vals.Set("client_id", c.ClientID)
		vals.Set("client_secret", c.ClientSecret)
	}
	for k, v := range params {
		vals.Set(k, v)
	}

	u := url.URL{
		Scheme:   "https",
		Path:     path,
		RawQuery: vals.Encode(),
	}
	return http.Get(u.String())
}
