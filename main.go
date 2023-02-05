package main

import (
	"eiffel65/steam"
	"encoding/json"
	"flag"
	"fmt"
	"log"
)

var (
	assetName   string
	steamAPIKey string
	wearTier    int
	listings    int
	statTrak    bool
	debug       bool
)

const (
	defaultWearTier     int    = 3
	defaultListingCount int    = 25
	defaultAssetName    string = "AK-47 | Case Hardened"
)

func init() {
	flag.StringVar(&assetName, "n", defaultAssetName, "the name of the Steam asset to query")
	flag.StringVar(&steamAPIKey, "k", "", "the user Steam Web API Key")
	flag.IntVar(&wearTier, "w", defaultWearTier, "what wear quality to query (1-5 Factory New to Battle-Scarred, default 3)")
	flag.IntVar(&listings, "l", defaultListingCount, "how many market listings, default 25")
	flag.BoolVar(&statTrak, "s", false, "whether to query items with StatTrak")
	flag.BoolVar(&debug, "d", false, "debug mode")
	flag.Parse()
}

func main() {
	if steamAPIKey == "" {
		log.Fatal("please specify an API Key")
	}

	if wearTier > 5 || wearTier < 1 {
		log.Fatal("please specify a wear tear between 1 and 5")
	}

	steamClient := steam.NewClient(steamAPIKey)

	assetList, err := steamClient.NewAsset(assetName, wearTier, listings, statTrak, debug)
	if err != nil {
		log.Fatalf("failed to get asset listings for %s", err)
	}

	if assetList == nil {
		log.Fatalf("no results for %s", assetName)
	}

	assetJSON, err := json.MarshalIndent(assetList, "", "\t")
	if err != nil {
		log.Fatalf("failed to marshal listing JSON: %s", err)
	}

	notableIDs := steam.CheckForRarity(*assetList)
	highlight := ""
	if len(notableIDs) > 0 {
		for id, asset := range notableIDs {
			highlight += fmt.Sprintf("\nHIGHLIGHT: %s SEED: %d FLOAT: %.4f PRICE: %s%s SCREENSHOT: %s",
				id, asset.Float.PaintSeed, asset.Float.FloatValue, asset.ListingTotalPrice, asset.ListingCurrency, asset.ScreenshotURL)
		}
	}

	fmt.Printf("%s\n\n%s", assetJSON, highlight)
}
