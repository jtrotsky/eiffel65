package steam

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
)

const (
	marketBaseURL       string = "https://steamcommunity.com"
	marketLanguage      string = "en_US"
	marketCountry       string = "us"
	marketCurrency      string = "1" // USD (https://github.com/SteamRE/SteamKit/blob/master/Resources/SteamLanguage/enums.steamd#L855).
	marketDataFormat    string = "json"
	marketStartingIndex int    = 0

	pathPriceOverview string = "market/priceoverview"
	pathMarketListing string = "market/listings"
)

// MarketListing is an item listed on the Steam market.
type MarketListing struct {
	PageSize   string `json:"pagesize,omitempty"`
	TotalCount int    `json:"total_count,omitempty"`
	Start      int    `json:"start,omitempty"`

	ListingInfo map[string]Listing                     `json:"listinginfo,omitempty"`
	Assets      map[string]map[string]map[string]Asset `json:"assets,omitempty"`
}

// Listing contains information specific to the market listing such as its price.
type Listing struct {
	ID    string      `json:"listingid,omitempty"`
	Price int         `json:"price,omitempty"`
	Fee   int         `json:"fee,omitempty"`
	Asset MarketAsset `json:"asset,omitempty"`
}

// MarketAsset is a summary of an asset listed on the market.
type MarketAsset struct {
	ID            string         `json:"id,omitempty"`
	Currency      int            `json:"currency,omitempty"`
	AppID         int            `json:"appid,omitempty"`
	ContextID     string         `json:"contextid,omitempty"`
	Amount        string         `json:"amount,omitempty"`
	MarketActions []MarketAction `json:"market_actions,omitempty"`
}

// MarketAction is an action that can be taken such as inspect the item in game.
type MarketAction struct {
	Link string `json:"link,omitempty"`
	Name string `json:"name,omitempty"`
}

// GetMarketListing returns info about an asset listed on the Steam market.
func (client *Client) GetMarketListing(encodedName string, listings int, debug bool) (*MarketListing, error) {
	marketListingURL, err := url.Parse(fmt.Sprintf("%s/%s/%s/%s/render", marketBaseURL, pathMarketListing, csgoAppID, encodedName))
	if err != nil {
		return nil, err
	}

	params := url.Values{}
	params.Add("start", strconv.Itoa(marketStartingIndex))
	params.Add("count", strconv.Itoa(listings))
	params.Add("currency", marketCurrency)
	params.Add("language", marketLanguage)
	params.Add("format", marketDataFormat)
	params.Add("appid", csgoAppID)
	marketListingURL.RawQuery = params.Encode()

	// DEBUG
	if debug {
		log.Println(marketListingURL)
	}

	response, err := http.DefaultClient.Get(marketListingURL.String())
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	marketListing := MarketListing{}
	err = json.NewDecoder(response.Body).Decode(&marketListing)
	if err != nil {
		return nil, err
	}

	return &marketListing, err
}

// GetPriceSummary returns basic market price statistics for a given asset.
func (asset *SimpleAsset) GetPriceSummary(debug bool) error {
	assetPriceURL, err := url.Parse(fmt.Sprintf("%s/%s", marketBaseURL, pathPriceOverview))
	if err != nil {
		return err
	}

	params := url.Values{}
	params.Add("format", marketDataFormat)
	params.Add("country", marketCountry)
	params.Add("currency", marketCurrency)
	params.Add("appid", csgoAppID)
	assetPriceURL.RawQuery = params.Encode()

	assetPriceURL.RawQuery += fmt.Sprintf("&market_hash_name=%s", asset.EncodedName)

	// DEBUG
	if debug {
		log.Println(assetPriceURL)
	}

	response, err := http.DefaultClient.Get(assetPriceURL.String())
	if err != nil {
		return err
	}
	defer response.Body.Close()

	err = json.NewDecoder(response.Body).Decode(&asset.MarketValue)
	if err != nil {
		return err
	}

	if marketCurrency == "1" {
		asset.MarketValue.Currency = "USD"
	} else {
		asset.MarketValue.Currency = "?"
	}

	return err
}
