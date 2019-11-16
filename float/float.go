package float

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

const (
	csgoFloatBaseURL string = "https://api.csgofloat.com/"
)

// AssetFloatPayload contains the payload of data on asset quality and appearance.
type AssetFloatPayload struct {
	ItemInfo AssetFloat `json:"iteminfo,omitempty"`
}

// AssetFloat contains readings of the asset quality and appearance.
type AssetFloat struct {
	AccountID  string    `json:"accountid,omitempty"`
	DefIndex   int       `json:"defindex,omitempty"`
	PaintIndex int       `json:"paintindex,omitempty"`
	Rarity     int       `json:"rarity,omitempty"`
	Quality    int       `json:"quality,omitempty"`
	PaintWear  int64     `json:"paintwear,omitempty"`
	PaintSeed  int       `json:"paintseed,omitempty"`
	CustomName string    `json:"customname,omitempty"`
	Stickers   []Sticker `json:"stickers,omitempty"`
	Inventory  int       `json:"inventory,omitempty"`
	Origin     int       `json:"origin,omitempty"`
	FloatValue float64   `json:"floatvalue,omitempty"`
	ItemID     int64     `json:"itemid_int,omitempty"`
	WeaponType string    `json:"weapon_type,omitempty"`
	ItemName   string    `json:"item_name,omitempty"`
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

// Get looks up the asset paint/design quality.
func Get(inspectURL string) (*AssetFloatPayload, string, error) {
	csgoFloatURL, err := url.Parse(fmt.Sprintf("%s?url=%s", csgoFloatBaseURL, inspectURL))
	if err != nil {
		return nil, csgoFloatURL.String(), err
	}

	response, err := http.DefaultClient.Get(csgoFloatURL.String())
	if err != nil {
		return nil, csgoFloatURL.String(), err
	}
	defer response.Body.Close()

	assetFloatPayload := AssetFloatPayload{}
	err = json.NewDecoder(response.Body).Decode(&assetFloatPayload)
	if err != nil {
		return nil, csgoFloatURL.String(), err
	}

	return &assetFloatPayload, csgoFloatURL.String(), nil
}
