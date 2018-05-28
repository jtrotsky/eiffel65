package steam

const (
	csgoAppID               string = "730" // Steam Application ID for CSGO.
	steamAPIBaseURL         string = "https://api.steampowered.com"
	steamImageCDNBaseURL    string = "https://steamcommunity-a.akamaihd.net/economy/image/"
	steamImageCDNBaseURLOld string = "https://cdn.steamcommunity.com/economy/image/"

	csgoFloatBaseURL string = "https://api.csgofloat.com:1738/"
)

// Client is the Steam client that contains config and authentication.
type Client struct {
	APIKey           string
	CSGOAppID        string
	CSGOFloatBaseURL string
	CDNBaseURL       string
	APIBaseURL       string
}

// NewClient initiates a Steam client.
func NewClient(apiKey string) *Client {
	return &Client{
		APIKey:           apiKey,
		CSGOAppID:        csgoAppID,
		CDNBaseURL:       steamImageCDNBaseURL,
		APIBaseURL:       steamAPIBaseURL,
		CSGOFloatBaseURL: csgoFloatBaseURL,
	}
}
