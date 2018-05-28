package steam

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
)

const (
	statTrak string = "StatTrak™"

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
	Assets []AssetValue `json:"assets,omitempty"`

	Success bool `json:"success"`
}

// Asset is an item on the Steam market.
type Asset struct {
	ID                string         `json:"id,omitempty"`
	ClassID           string         `json:"classid,omitempty"`
	ContextID         string         `json:"contextid,omitempty"`
	InstanceID        string         `json:"instanceid,omitempty"`
	Name              string         `json:"name,omitempty"`
	Type              string         `json:"type,omitempty"`
	IconURL           string         `json:"icon_url,omitempty"`       // Prefaced by CDN URL.
	IconURLLarge      string         `json:"icon_url_large,omitempty"` // Prefaced by CDN URL.
	IconDragURL       string         `json:"icon_drag_url,omitempty"`  // Prefaced by CDN URL.
	MarketHashName    string         `json:"market_hash_name,omitempty"`
	MarketName        string         `json:"market_name,omitempty"`
	NameColor         string         `json:"name_color,omitempty"`
	BGColor           string         `json:"background_color,omitempty"`
	Tradable          int            `json:"tradable,omitempty"`
	Marketable        int            `json:"marketable,omitempty"`
	Commodity         int            `json:"commodity,omitempty"`
	TradeRestrict     int            `json:"market_tradeable_restriction,omitempty"`
	FraudWarnings     string         `json:"fraudwarnings,omitempty"`
	OwnerDescriptions string         `json:"owner_descriptions,omitempty"`
	MarketActions     []MarketAction `json:"market_actions,omitempty"`
	Descriptions      []Description  `json:"descriptions,omitempty"`
	Tags              []Tag          `json:"tags,omitempty"`
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
	ID          string       `json:"id,omitempty"`
	ClassID     string       `json:"class_id,omitempty"`
	ContextID   string       `json:"context_id,omitempty"`
	InstanceID  string       `json:"instance_id,omitempty"`
	Name        string       `json:"name,omitempty"`
	EncodedName string       `json:"encoded_name,omitempty"`
	IconURL     string       `json:"icon_url,omitempty"`
	InspectURL  string       `json:"inspect_url,omitempty"`
	Type        AssetType    `json:"type,omitempty"`
	MarketValue AssetValue   `json:"market_value,omitempty"`
	Quality     AssetQuality `json:"quality,omitempty"`
	Float       AssetFloat   `json:"float,omitempty"`
}

// AssetQuality is the weapon condition and rarity.
type AssetQuality struct {
	Wear AssetWear `json:"wear,omitempty"`
	Type string    `json:"type,omitempty"`
}

// AssetValue contains asset price statistics.
type AssetValue struct {
	Currency    string `json:"currency,omitempty"`
	LowestPrice string `json:"lowest_price,omitempty"`
	MedianPrice string `json:"median_price,omitempty"`
	Volume      string `json:"volume,omitempty"`
}

// AssetFloatPayload contains the payload of data on asset quality and appearance.
type AssetFloatPayload struct {
	ItemInfo AssetFloat `json:"iteminfo,omitempty"`
}

// AssetFloat contains readings of the asset quality and appearance.
type AssetFloat struct {
	AccountID string `json:"accountid,omitempty"`
	// "itemid": {
	// 	"low": -1829587929,
	// 	"high": 1,
	// 	"unsigned": true
	// },
	DefIndex   int   `json:"defindex,omitempty"`
	PaintIndex int   `json:"paintindex,omitempty"`
	Rarity     int   `json:"rarity,omitempty"`
	Quality    int   `json:"quality,omitempty"`
	PaintWear  int64 `json:"paintwear,omitempty"`
	PaintSeed  int64 `json:"paintseed,omitempty"`
	// "killeaterscoretype": 0,
	// "killeatervalue": 0,
	// "customname": null,
	Stickers  []Sticker `json:"stickers,omitempty"`
	Inventory int       `json:"inventory,omitempty"`
	Origin    int       `json:"origin,omitempty"`
	// "questid": null,
	// "dropreason": null,
	FloatValue float64 `json:"floatvalue,omitempty"`
	ItemID     int64   `json:"itemid_int,omitempty"`
	// "s": "0",
	// "a": "6760346663",
	// "d": "30614827701953021",
	// "m": "625254122282020305",
	// "imageurl": "http://media.steampowered.com/apps/730/icons/econ/default_generated/weapon_ak47_cu_ak47_cobra_light_large.7494bfdf4855fd4e6a2dbd983ed0a243c80ef830.png",
	// "min": 0.1,
	// "max": 0.7,
	WeaponType string `json:"weapon_type,omitempty"`
	ItemName   string `json:"item_name,omitempty"`
}

// Sticker is an asset that can be attached to CSGO weapons.
type Sticker struct {
	Slot      int     `json:"slot,omitempty"`
	StickerID int64   `json:"sticker_id,omitempty"`
	Wear      float64 `json:"wear,omitempty"`
	Scale     float64 `json:"scale,omitempty"`
	Rotation  float64 `json:"rotation,omitempty"`
	CodeName  string  `json:"codename,omitempty"`
	Name      string  `json:"name,omitempty"`
}

// NewAsset creates an asset instance.
func (client *Client) NewAsset(name string, wearTier int, isStatTrak bool) *SimpleAsset {
	wear := getWearTierName(wearTier)
	marketName := formatMarketName(name, wear, isStatTrak)

	simpleAsset := SimpleAsset{
		Name:        marketName,
		EncodedName: url.PathEscape(marketName),
		Type:        weaponAsset,
		Quality: AssetQuality{
			Wear: wear,
		},
	}

	marketListing, err := client.GetMarketListing(simpleAsset.EncodedName)
	if err != nil {
		log.Fatalf("failed get market listing: %s", err)
	}

	// Get unknown ClassID from market listing.
	classID := ""
	for k := range marketListing.Assets[client.CSGOAppID]["2"] {
		classID = k
	}

	// assetInfo, err := client.GetAsset(classID)
	// if err != nil {
	// 	log.Fatalf("failed get asset info: %s", err)
	// }

	err = simpleAsset.GetPriceSummary()
	if err != nil {
		log.Fatalf("failed get price summary: %s", err)
	}

	simpleAsset.ID = marketListing.Assets[client.CSGOAppID]["2"][classID].ID
	simpleAsset.ClassID = classID
	simpleAsset.ContextID = marketListing.Assets[client.CSGOAppID]["2"][classID].ContextID
	simpleAsset.InstanceID = marketListing.Assets[client.CSGOAppID]["2"][classID].InstanceID
	simpleAsset.Quality.Type = marketListing.Assets[client.CSGOAppID]["2"][classID].Type

	for _, action := range marketListing.Assets[client.CSGOAppID]["2"][classID].MarketActions {
		if action.Name == "Inspect in Game..." {
			simpleAsset.InspectURL = parseInspectURL(simpleAsset.ID, action.Link)
		}
	}

	if marketListing.Assets[client.CSGOAppID]["2"][classID].IconURLLarge != "" {
		simpleAsset.IconURL += fmt.Sprintf(client.CDNBaseURL + marketListing.Assets[client.CSGOAppID]["2"][classID].IconURLLarge)
	}

	if simpleAsset.InspectURL != "" {
		assetFloat, err := client.getAssetFloat(simpleAsset.InspectURL)
		if err != nil {
			log.Fatalf("failed get price summary: %s", err)
		}

		simpleAsset.Float = assetFloat.ItemInfo
	}

	return &simpleAsset
}

// getAssetPrices returns a list of asset prices for a given AppId.
func (client *Client) getAssetPrices() (*AssetPayload, error) {
	assetPricesURL, err := url.Parse(fmt.Sprintf("%s/%s", steamAPIBaseURL, pathAssetPrices))
	if err != nil {
		return nil, err
	}

	params := url.Values{}
	params.Add("currency", marketCurrency)
	params.Add("language", marketLanguage)
	params.Add("appid", csgoAppID)
	params.Add("key", client.APIKey)
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
func (client *Client) GetAsset(classID string) (*Asset, error) {
	assetInfoURL, err := url.Parse(fmt.Sprintf("%s/%s", steamAPIBaseURL, pathAssetInfo))
	if err != nil {
		return nil, err
	}

	params := url.Values{}
	// params.Add("instanceid0", instanceid) // Optional
	params.Add("class_count", "1") // Number of classes specified.
	params.Add("classid0", classID)
	params.Add("appid", csgoAppID)
	params.Add("key", client.APIKey)
	assetInfoURL.RawQuery = params.Encode()

	log.Println(assetInfoURL)

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

// getWearTierName identifies the wear quality category of an asset.
func getWearTierName(wearTier int) AssetWear {
	var wear AssetWear
	switch wearTier {
	case 1:
		wear = factoryNew
	case 2:
		wear = minimalWear
	case 3:
		wear = fieldTested
	case 4:
		wear = wellWorn
	case 5:
		wear = battleScared
	}
	return wear
}

// formatMarketName creates the Steam market-searchable name for an asset.
func formatMarketName(baseName string, wear AssetWear, isStatTrak bool) string {
	// StatTrak™ AK-47 | Case Hardened (Field-Tested)
	marketName := ""
	if isStatTrak {
		marketName = statTrak + " "
	}

	marketName += baseName + " " + "(" + string(wear) + ")"
	return marketName
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

// getAssetFloat looks up the asset paint/design quality.
func (client *Client) getAssetFloat(inspectURL string) (*AssetFloatPayload, error) {
	csgoFloatURL, err := url.Parse(fmt.Sprintf("%s?url=%s", client.CSGOFloatBaseURL, inspectURL))
	if err != nil {
		return nil, err
	}

	log.Println(csgoFloatURL)

	response, err := http.DefaultClient.Get(csgoFloatURL.String())
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	assetFloatPayload := AssetFloatPayload{}
	err = json.NewDecoder(response.Body).Decode(&assetFloatPayload)
	if err != nil {
		return nil, err
	}

	return &assetFloatPayload, nil
}

// parseInspectURL takes a raw inspect URL and converts it to one that can be
// looked up with CSGOFloat https://github.com/Step7750/CSGOFloat
func parseInspectURL(assetID, rawInspectURL string) string {
	inspectURL := strings.Split(rawInspectURL, "/")
	inspectURLTrimmed := strings.TrimPrefix(inspectURL[5], "+csgo_econ_action_preview%20")
	inspectURLJoined := "steam://rungame/730/76561202255233023/+csgo_econ_action_preview%20"

	d := strings.Index(inspectURLTrimmed, "D")
	dID := inspectURLTrimmed[d:]

	inspectURLTrimmed = strings.TrimSuffix(inspectURLTrimmed, inspectURLTrimmed[d:])

	a := strings.Index(inspectURLTrimmed, "A")
	aID := inspectURLTrimmed[a:]

	inspectURLTrimmed = strings.TrimSuffix(inspectURLTrimmed, inspectURLTrimmed[a:])

	aID = "A" + assetID

	m := strings.Index(inspectURLTrimmed, "M")
	mID := inspectURLTrimmed[m:]

	inspectURLJoined += mID + aID + dID

	return inspectURLJoined
}

// ListAssets prints assets in a list.
// func (client *Client) ListAssets(apiKey string) error {
// 	assetPayload, err := client.getAssetPrices()
// 	if err != nil {
// 		return err
// 	}

// 	assetList := []SimpleAsset{}
// 	count := 0
// 	for _, assetValue := range assetPayload.Result.Assets {
// 		log.Printf("getting: %s", assetValue.ClassID)

// 		asset, err := client.GetAsset(assetValue.ClassID)
// 		if err != nil {
// 			log.Println(err)
// 		}

// 		simpleAsset, err := asset.Transform()
// 		if err != nil {
// 			log.Println(err)
// 		}

// 		// if simpleAsset.Type != "weapon" {
// 		// 	log.Println("type not weapon")

// 		// 	continue
// 		// }

// 		count++

// 		assetList = append(assetList, *simpleAsset)
// 		if count == 10 {
// 			break
// 		}
// 	}

// 	assetJSON, err := json.MarshalIndent(assetList, "", "\t")
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	fmt.Printf("%s\n", assetJSON)

// 	return err
// }
