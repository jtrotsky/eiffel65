package steam

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
)

const (
	marketBaseURL    string = "https://steamcommunity.com"
	marketLanguage   string = "en_US"
	marketCountry    string = "us"
	marketCurrency   string = "1" // USD (https://github.com/SteamRE/SteamKit/blob/master/Resources/SteamLanguage/enums.steamd#L855).
	marketDataFormat string = "json"

	pathPriceOverview string = "market/priceoverview"
	pathMarketListing string = "market/listings"
)

// MarketListing is an item listed on the Steam market.
type MarketListing struct {
	PageSize   string `json:"pagesize,omitempty"`
	TotalCount int    `json:"total_count,omitempty"`
	Start      int    `json:"start,omitempty"`

	ListingInfo map[string]Listing `json:"listinginfo,omitempty"`
	Assets      map[string]string  `json:"assets,omitempty"`
}

// Listing contains information specific to the market listing such as its price.
type Listing struct {
	ID    string `json:"listingid,omitempty"`
	Price int64  `json:"price,omitempty"`
	Fee   int64  `json:"fee,omitempty"`

	// "publisher_fee_app": 730,
	// "publisher_fee_percent": "0.100000001490116119",
	// "currencyid": 2001,
	// "steam_fee": 173,
	// "publisher_fee": 347,
	// "converted_price": 2988,
	// "converted_fee": 447,
	// "converted_currencyid": 2003,
	// "converted_steam_fee": 149,
	// "converted_publisher_fee": 298,
	// "converted_price_per_unit": 2988,
	// "converted_fee_per_unit": 447,
	// "converted_steam_fee_per_unit": 149,
	// "converted_publisher_fee_per_unit": 298,

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
func GetMarketListing(encodedName string) (*MarketListing, error) {
	marketListingURL, err := url.Parse(fmt.Sprintf("%s/%s/%s/%s/render", marketBaseURL, pathMarketListing, steamAppID, encodedName))
	if err != nil {
		return nil, err
	}

	params := url.Values{}
	params.Add("start", "0") // consts
	params.Add("count", "1") // consts
	params.Add("currency", marketCurrency)
	params.Add("language", marketLanguage)
	params.Add("format", marketDataFormat)
	params.Add("appid", steamAppID)
	marketListingURL.RawQuery = params.Encode()

	// DEBUG
	log.Println(marketListingURL)

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

	log.Printf("%s", marketListing.Assets)

	return &marketListing, err
}

// GetPriceSummary returns basic market price statistics for a given asset.
func (asset *SimpleAsset) GetPriceSummary() error {
	assetPriceURL, err := url.Parse(fmt.Sprintf("%s/%s", marketBaseURL, pathPriceOverview))
	if err != nil {
		return err
	}

	params := url.Values{}
	params.Add("format", marketDataFormat)
	params.Add("country", marketCountry)
	params.Add("currency", marketCurrency)
	params.Add("appid", steamAppID)
	assetPriceURL.RawQuery = params.Encode()

	assetPriceURL.RawQuery += fmt.Sprintf("&market_hash_name=%s", url.PathEscape(asset.EncodedName))

	// DEBUG
	log.Println(assetPriceURL)

	response, err := http.DefaultClient.Get(assetPriceURL.String())
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.ContentLength <= 2 {
		return errors.New("failed to get asset price summary")
	}

	err = json.NewDecoder(response.Body).Decode(&asset.Price.Summary)
	if err != nil {
		return err
	}

	if marketCurrency == "1" {
		asset.Price.Summary.Currency = "USD"
	} else {
		asset.Price.Summary.Currency = "?"
	}

	return err
}
