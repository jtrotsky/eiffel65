package steam

const (
	csgoAppID               string = "730" // Steam Application ID for CSGO.
	steamAPIBaseURL         string = "https://api.steampowered.com"
	steamImageCDNBaseURL    string = "https://steamcommunity-a.akamaihd.net/economy/image/"
	steamImageCDNBaseURLOld string = "https://cdn.steamcommunity.com/economy/image/"
)

// Client is the Steam client that contains config and authentication.
type Client struct {
	APIKey     string
	CSGOAppID  string
	CDNBaseURL string
	APIBaseURL string
	Debug      bool
}

// NewClient initiates a Steam client.
func NewClient(apiKey string, debug bool) *Client {
	return &Client{
		APIKey:     apiKey,
		CSGOAppID:  csgoAppID,
		CDNBaseURL: steamImageCDNBaseURL,
		APIBaseURL: steamAPIBaseURL,
		Debug:      debug,
	}
}
