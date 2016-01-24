package foursquare

const (
	DefaultScheme = "https"

	// Endpoints for the OAuth 2.0 flow
	OAuthRoot     = "foursquare.com/oauth2"
	AuthEndpoint  = "/authenticate"
	TokenEndpoint = "/access_token"

	// API endpoints
	APIRoot    = "api.foursquare.com"
	APIVersion = "v2"

	// See: https://developer.foursquare.com/overview/versioning
	// The `v` param: version in YYYYMMDD format
	DateVersion = "20160124"
	// The `m` param: Foursquare or Swarm mode
	FoursquareMode = "foursquare"
	SwarmMode      = "iota"
)
