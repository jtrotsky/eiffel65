package steam

const (
	CSGOApplicationID       int    = 730 // Steam Application ID for CSGO.
	steamAPIBaseURL         string = "https://api.steampowered.com"
	steamImageCDNBaseURL    string = "https://steamcommunity-a.akamaihd.net/economy/image/"
	steamImageCDNBaseURLOld string = "https://cdn.steamcommunity.com/economy/image/"
)

// Client is the Steam client that contains config and authentication.
type Client struct {
	APIKey            string
	CSGOApplicationID int
	CDNBaseURL        string
	APIBaseURL        string
}

// NewClient initiates a Steam client.
func NewClient(apiKey string) *Client {
	return &Client{
		APIKey:            apiKey,
		CSGOApplicationID: CSGOApplicationID,
		CDNBaseURL:        steamImageCDNBaseURL,
		APIBaseURL:        steamAPIBaseURL,
	}
}
