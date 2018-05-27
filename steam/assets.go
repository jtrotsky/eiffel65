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
	isStatTrak bool = false

	factoryNew   AssetWear = "Factory New"
	minimalWear  AssetWear = "Minimal Wear"
	fieldTested  AssetWear = "Field-Tested"
	wellWorn     AssetWear = "Well-Worn"
	battleScared AssetWear = "Battle-Scarred"

	badgeAsset       AssetType = "badge"
	caseAsset        AssetType = "case"
	collectableAsset AssetType = "collectable"
	glovesAsset      AssetType = "gloves"
	keyAsset         AssetType = "key"
	tagAsset         AssetType = "tag"
	toolAsset        AssetType = "tool"
	weaponAsset      AssetType = "weapon"

	pathAssetInfo   string = "ISteamEconomy/GetAssetClassInfo/v0001"
	pathAssetPrices string = "ISteamEconomy/GetAssetPrices/v1"
)

// AssetWear is how assets are categorised by quality based on their
// appearance.
type AssetWear string

// AssetType is a type of asset on the Steam market.
type AssetType string

// AssetPayload is to make docoding API payloads easier.
type AssetPayload struct {
	Result AssetPayloadResult `json:"result"`
}

// AssetPayloadResult is the second layer of an asset payload.
type AssetPayloadResult struct {
	Assets []AssetPrice `json:"assets,omitempty"`

	Success bool `json:"success"`
}

// Asset is an item on the Steam market.
type Asset struct {
	Name              string                 `json:"name,omitempty"`
	Type              string                 `json:"type,omitempty"`
	IconURL           string                 `json:"icon_url,omitempty"` // Prefaced by http://cdn.steamcommunity.com/economy/image/
	IconURLLarge      string                 `json:"icon_url_large,omitempty"`
	IconDragURL       string                 `json:"icon_drag_url,omitempty"`
	MarketHashName    string                 `json:"market_hash_name,omitempty"`
	MarketName        string                 `json:"market_name,omitempty"`
	NameColor         string                 `json:"name_color,omitempty"`
	BGColor           string                 `json:"background_color,omitempty"`
	Tradable          string                 `json:"tradable,omitempty"`
	Marketable        string                 `json:"marketable,omitempty"`
	Commodity         string                 `json:"commodity,omitempty"`
	TradeRestrict     string                 `json:"market_tradeable_restriction,omitempty"`
	FraudWarnings     string                 `json:"fraudwarnings,omitempty"`
	Descriptions      map[string]Description `json:"descriptions,omitempty"`
	OwnerDescriptions string                 `json:"owner_descriptions,omitempty"`
	Tags              map[string]Tag         `json:"tags,omitempty"`
	ClassID           string                 `json:"classid,omitempty"`
}

// Description contains asset description data
type Description struct {
	Type    string `json:"type"`
	Value   string `json:"value"`
	Color   string `json:"color,omitempty"`
	AppData string `json:"appdata"`
}

// Tag contains asset tag data
type Tag struct {
	InternalName string `json:"internal_name"`
	Name         string `json:"name"`
	Category     string `json:"category"`
	Color        string `json:"color,omitempty"`
	CategoryName string `json:"category_name"`
}

// SimpleAsset is a simple version of Asset.
type SimpleAsset struct {
	Name        string       `json:"name,omitempty"`
	EncodedName string       `json:"encoded_name,omitempty"`
	Type        AssetType    `json:"type,omitempty"`
	Price       AssetPrice   `json:"price,omitempty"`
	Quality     AssetQuality `json:"quality,omitempty"`
}

// AssetQuality is the weapon condition and rarity.
type AssetQuality struct {
	Wear AssetWear `json:"wear,omitempty"`
}

// AssetPrice contains asset price statistics.
type AssetPrice struct {
	Name    string            `json:"name,omitempty"`
	Date    string            `json:"date,omitempty"`
	ClassID string            `json:"classid,omitempty"`
	Summary AssetPriceSummary `json:"summary,omitempty"`
}

// AssetPriceSummary contains basic market price statistics about an asset.
type AssetPriceSummary struct {
	Currency    string `json:"currency,omitempty"`
	LowestPrice string `json:"lowest_price,omitempty"`
	MedianPrice string `json:"median_price,omitempty"`
	Volume      string `json:"volume,omitempty"`
}

// NewAsset creates an asset instance.
func NewAsset(name string) *SimpleAsset {
	// Name: "StatTrak\u2122 AK-47 | Blue Laminate (Factory New)",
	return &SimpleAsset{
		Name:        name,
		EncodedName: url.PathEscape(name),
	}
}

// ListAssets prints assets in a list.
func ListAssets(apiKey string) error {
	assetPayload, err := getAssetPrices(apiKey)
	if err != nil {
		return err
	}

	assetList := []SimpleAsset{}
	count := 0
	for _, assetPrice := range assetPayload.Result.Assets {
		log.Printf("getting: %s", assetPrice.ClassID)

		asset, err := GetAsset(apiKey, assetPrice.ClassID)
		if err != nil {
			log.Println(err)
		}

		simpleAsset, err := asset.Transform()
		if err != nil {
			log.Println(err)
		}

		// if simpleAsset.Type != "weapon" {
		// 	log.Println("type not weapon")

		// 	continue
		// }

		count++

		assetList = append(assetList, *simpleAsset)
		if count == 10 {
			break
		}
	}

	assetJSON, err := json.MarshalIndent(assetList, "", "\t")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%s\n", assetJSON)

	return err
}

// getAssetPrices returns a list of asset prices for a given AppId.
func getAssetPrices(steamAPIKey string) (*AssetPayload, error) {
	assetPricesURL, err := url.Parse(fmt.Sprintf("%s/%s", steamAPIBaseURL, pathAssetPrices))
	if err != nil {
		return nil, err
	}

	params := url.Values{}
	params.Add("currency", marketCurrency)
	params.Add("language", marketLanguage)
	params.Add("appid", steamAppID)
	params.Add("key", steamAPIKey)
	assetPricesURL.RawQuery = params.Encode()

	response, err := http.DefaultClient.Get(assetPricesURL.String())
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	assetPayload := AssetPayload{}
	err = json.NewDecoder(response.Body).Decode(&assetPayload)
	if err != nil {
		return nil, err
	}

	if assetPayload.Result.Success != true {
		return nil, errors.New("failed getting asset price list")
	}

	return &assetPayload, err
}

// GetAsset returns information about an individual item from the Steam market.
func GetAsset(apiKey, classID string) (*Asset, error) {
	assetInfoURL, err := url.Parse(fmt.Sprintf("%s/%s", steamAPIBaseURL, pathAssetInfo))
	if err != nil {
		return nil, err
	}

	params := url.Values{}
	// params.Add("instanceid0", instanceid) // Optional
	params.Add("class_count", "1") // Number of classes specified.
	params.Add("classid0", classID)
	params.Add("appid", steamAppID)
	params.Add("key", apiKey)
	assetInfoURL.RawQuery = params.Encode()

	response, err := http.DefaultClient.Get(assetInfoURL.String())
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	type Payload struct {
		Result map[string]json.RawMessage `json:"result,omitempty"`
	}
	payload := Payload{}

	err = json.NewDecoder(response.Body).Decode(&payload)
	if err != nil {
		return nil, err
	}

	asset := Asset{}
	for k, v := range payload.Result {
		// TODO check if success is false.
		// if k == "success" {
		// 	if v.(bool) == false {
		// 		return nil, errors.New("failed to get asset")
		// 	}
		// }

		if k == classID {
			err := json.Unmarshal(v, &asset)
			if err != nil {
				return nil, errors.New("failed to unmarshal asset")
			}

			return &asset, nil
		}
	}

	return nil, err
}

// Transform converts a raw Steam asset into a simplified one.
func (asset *Asset) Transform() (*SimpleAsset, error) {
	assetSimple := SimpleAsset{}

	assetSimple.Name = asset.MarketName
	assetSimple.EncodedName = asset.MarketHashName

	switch asset.Type {
	case "Base Grade Collectible":
		assetSimple.Type = collectableAsset
	case "Base Grade Container":
		assetSimple.Type = caseAsset
	case "Base Grade Key":
		assetSimple.Type = keyAsset
	case "Base Grade Tag":
		assetSimple.Type = tagAsset
	case "Base Grade Tool":
		assetSimple.Type = toolAsset
	default:
		fmt.Println(asset.Type)
	}

	err := assetSimple.GetPriceSummary()
	if err != nil {
		return nil, err
	}

	return &assetSimple, err
}
