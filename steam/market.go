package steam

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
)

//https://steamcommunity.com/market/search/render/
//?sort_column=default
//&sort_dir=desc
//&appid=730
//&norender=1
//&count=25
//&query=AK-47%20%7C%20Case%20Hardened%20%28Field-Tested%29

const (
	marketBaseURL    string = "https://steamcommunity.com"
	marketSearchPath string = "/market/search/render/"

	// pathPriceOverview string = "market/priceoverview"
	// pathMarketListing string = "market/listings"
)

// MarketSearch is a search for an item on the Steam market.
type MarketSearch struct {
	Success    bool `json:"success,omitempty"`
	PageSize   int  `json:"pagesize,omitempty"`
	TotalCount int  `json:"total_count,omitempty"`
	Start      int  `json:"start,omitempty"`

	SearchData SearchData `json:"searchdata,omitempty"`
	Results    []Listing
}

// SearchData summarises the requested market search and the results found.
type SearchData struct {
	Query              string `json:"query,omitempty"`
	SearchDescriptions bool   `json:"search_descriptions,omitempty"`
	TotalCount         int    `json:"total_count,omitempty"`
	PageSize           int    `json:"pagesize,omitempty"`
	Prefix             string `json:"prefix,omitempty"`
	ClassPrefix        string `json:"class_prefix,omitempty"`
}

// Listing represents an item for sale on the market including its quality and price.
type Listing struct {
	Name             string           `json:"name,omitempty"`
	HashName         string           `json:"hash_name,omitempty"`
	SellListings     int              `json:"sell_listings,omitempty"`
	SellPrice        int              `json:"sell_price,omitempty"`
	SellPriceText    string           `json:"sell_price_text,omitempty"`
	AppIcon          string           `json:"app_icon,omitempty"`
	AppName          string           `json:"app_name,omitempty"`
	AssetDescription AssetDescription `json:"asset_description,omitempty"`
	SalePriceText    string           `json:"sale_price_text,omitempty"`
}

// AssetDescription contains details about an item.
type AssetDescription struct {
	AppID            int    `json:"appid,omitempty"`
	ClassID          string `json:"classid,omitempty"`
	InstanceID       string `json:"instanceid,omitempty"`
	BackgroundColour string `json:"appid,omitempty"`
	IconURL          string `json:"background_color,omitempty"`
	Tradable         bool   `json:"tradable,omitempty"`
	Name             string `json:"name,omitempty"`
	NameColour       string `json:"name_color,omitempty"`
	Type             string `json:"type,omitempty"`
	MarketName       string `json:"market_name,omitempty"`
	MarketHashName   string `json:"market_hash_name,omitempty"`
	Commodity        bool   `json:"commodity,omitempty"`
}

// Search returns listings for an asset on the Steam market.
func (client *Client) Search(encodedItemName string, searchCount int, debug bool) (*MarketSearch, error) {
	marketSearchURL, err := url.Parse(fmt.Sprintf("%s%s", marketBaseURL, marketSearchPath))
	if err != nil {
		return nil, err
	}

	params := url.Values{}
	params.Add("appid", strconv.Itoa(CSGOApplicationID))
	// params.Add("sort_column", "default")
	// params.Add("sort_dir", "desc")
	params.Add("norender", strconv.Itoa(1)) // get results as JSON not HTML
	params.Add("count", strconv.Itoa(searchCount))
	params.Add("query", encodedItemName)
	marketSearchURL.RawQuery = params.Encode()

	// DEBUG
	if debug {
		log.Println(marketSearchURL)
	}

	response, err := http.DefaultClient.Get(marketSearchURL.String())
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	marketSearch := MarketSearch{}
	err = json.NewDecoder(response.Body).Decode(&marketSearch)
	if err != nil {
		return nil, err
	}

	return &marketSearch, err
}
