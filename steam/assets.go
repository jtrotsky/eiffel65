package steam

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/jtrotsky/eiffel65/float"
	"github.com/jtrotsky/eiffel65/image"
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

// Mostly from https://steamcommunity.com/sharedfiles/filedetails/?id=380042859
// Examples
// https://files.opskins.media/file/opskins-patternindex/7_44_29.jpg
// https://files.opskins.media/file/opskins-patternindex/7_44_151.jpg
// https://files.opskins.media/file/opskins-patternindex/7_44_179.jpg
// https://files.opskins.media/file/opskins-patternindex/7_44_321.jpg
//
// https://files.opskins.media/file/opskins-patternindex/7_44_464.jpg
var (
	ak47RarePaintSeeds     = []int{29, 151, 179, 321, 464, 561, 661, 670, 760, 955}
	falchionRarePaintSeeds = []int{4, 10, 11, 13, 14, 20, 25, 27, 29, 30, 32, 34,
		38, 42, 46, 55, 56, 58, 61, 67, 73, 74, 79, 82, 91, 92, 98, 103, 106, 109,
		112, 115, 116, 126, 128, 129, 130, 136, 137, 138, 139, 144, 146, 147, 148,
		149, 151, 152, 155, 157, 166, 168, 169, 170, 171, 175, 176, 177, 179, 180, 182,
		187, 188, 189, 191, 194, 199, 202, 203, 205, 207, 208, 210, 212, 213, 214,
		216, 217, 222, 225, 226, 228, 230, 231, 233, 235, 236, 237, 238, 239, 241,
		243, 244, 245, 246, 248, 251, 302, 494, 764, 811, 917}
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
	FraudWarnings     []string       `json:"fraudwarnings,omitempty"`
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
	ID string `json:"id,omitempty"`
	// ClassID string `json:"class_id,omitempty"`
	// ContextID         string    `json:"context_id,omitempty"`
	// InstanceID        string    `json:"instance_id,omitempty"`
	Name              string    `json:"name,omitempty"`
	EncodedName       string    `json:"encoded_name,omitempty"`
	IconURL           string    `json:"icon_url,omitempty"`
	InspectURL        string    `json:"inspect_url,omitempty"`
	ScreenshotURL     string    `json:"screenshot_url,omitempty"`
	ListingID         string    `json:"listing_id,omitempty"`
	ListingCurrency   string    `json:"listing_currency,omitempty"`
	ListingPrice      string    `json:"listing_price,omitempty"`
	ListingFee        string    `json:"listing_fee,omitempty"`
	ListingTotalPrice string    `json:"listing_total_price,omitempty"`
	Type              AssetType `json:"type,omitempty"`
	//MarketValue       AssetValue       `json:"market_value,omitempty"`
	Quality AssetQuality     `json:"quality,omitempty"`
	Float   float.AssetFloat `json:"float,omitempty"`
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

// NewAsset creates an asset instance.
func (client *Client) NewAsset(name string, wearTier, listings int, isStatTrak, debug bool) (*[]SimpleAsset, error) {
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

	// Returns a page of commmunity market listings for the given asset.
	marketListing, err := client.GetMarketListing(simpleAsset.EncodedName, listings, debug)
	if err != nil {
		return nil, err
	}

	if len(marketListing.Assets) == 0 {
		return nil, err
	}

	if len(marketListing.ListingInfo) == 0 {
		return nil, err
	}

	// The ClassID is unique ID of each listing, which we do not know until
	// it is returned in the listing summary.
	// classID := ""

	// Create an empty simple asset for each listing.
	assetListing := SimpleAsset{}

	// The list of simple asset listings to return at the end.
	simpleAssetList := []SimpleAsset{}

	// Loop through each asset listing, the key being the ClassID.
	for cID, listing := range marketListing.Assets[client.CSGOAppID]["2"] {
		// classID = cID

		// Fill out the basic asset info that is the same for each listing.
		assetListing = simpleAsset

		assetListing.ID = cID
		// assetListing.ID = listing.ID
		// assetListing.ContextID = listing.ContextID
		// assetListing.InstanceID = listing.InstanceID
		assetListing.Quality.Type = listing.Type

		// Pricing from listing
		//assetListing.ListingPrice = marketListing.ListingInfo[].Listing

		for _, action := range listing.MarketActions {
			if action.Name == "Inspect in Game..." {
				assetListing.InspectURL = parseInspectURL(assetListing.ID, action.Link)
			}
		}

		if assetListing.InspectURL != "" {
			assetFloat, floatURL, err := float.Get(assetListing.InspectURL)
			if err != nil {
				log.Printf("failed get price summary: %s", err)
				break
			}

			if debug {
				log.Println(floatURL)
			}

			assetListing.Float = assetFloat.ItemInfo

			screenshotURL, err := image.BuildURL(assetListing.Float.DefIndex, assetListing.Float.PaintIndex, assetListing.Float.PaintSeed, assetListing.InspectURL)
			if err != nil {
				log.Printf("failed to get screenshot: %s", err)
			}

			if debug {
				log.Println(screenshotURL)
			}

			assetListing.ScreenshotURL = screenshotURL
		}

		for listingID, listing := range marketListing.ListingInfo {
			if listing.Asset.ID == assetListing.ID {
				assetListing.ListingID = listingID

				listingPriceFloat := float64(listing.Price) / 100
				listingFeeFloat := float64(listing.Fee) / 100

				assetListing.ListingPrice = strconv.FormatFloat(listingPriceFloat, 'f', 2, 64)
				assetListing.ListingFee = strconv.FormatFloat(listingFeeFloat, 'f', 2, 64)
				assetListing.ListingTotalPrice = strconv.FormatFloat(listingPriceFloat+listingFeeFloat, 'f', 2, 64)

				assetListing.ListingCurrency = "USD"
			}
		}

		simpleAssetList = append(simpleAssetList, assetListing)
	}

	return &simpleAssetList, err
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
func (asset *Asset) Transform(debug bool) (*SimpleAsset, error) {
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

	return &assetSimple, nil
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

// CheckForRarity loops through floats for market listings and highlights any
// standout values.
func CheckForRarity(assetList []SimpleAsset) map[string]SimpleAsset {
	notableListings := map[string]SimpleAsset{}
	for _, asset := range assetList {
		if rarePaintSeed(asset.Float.DefIndex, asset.Float.PaintSeed) {
			notableListings[asset.ID] = asset
		}
	}
	return notableListings
}

// rarePaintSeed checks whether a seed exists in a list of rare ones.
func rarePaintSeed(defIndex, seed int) bool {
	switch defIndex {
	case 512:
		for _, rareSeed := range falchionRarePaintSeeds {
			if seed == rareSeed {
				return true
			}
		}
		return false
	case 7:
		for _, rareSeed := range ak47RarePaintSeeds {
			if seed == rareSeed {
				return true
			}
		}
		return false
	}
	return false
}
