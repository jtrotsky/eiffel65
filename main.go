package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"

	"github.com/jtrotsky/eiffel65/steam"
)

const (
	defaultAssetName string = "AK-47 | Case Hardened"
	defaultWearTier  int    = 3
)

var (
	assetName   string
	steamAPIKey string
	wearTier    int
	statTrak    bool
)

func init() {
	flag.StringVar(&assetName, "n", defaultAssetName, "the name of the Steam asset to query")
	flag.StringVar(&steamAPIKey, "k", "", "the user Steam Web API Key")
	flag.IntVar(&wearTier, "w", defaultWearTier, "what wear quality to query")
	flag.BoolVar(&statTrak, "s", false, "whether to query items with StatTrak")
	flag.Parse()
}

func main() {
	if steamAPIKey == "" {
		log.Fatal("please specify an API Key")
	}

	if wearTier > 5 || wearTier < 1 {
		log.Fatal("please specify a wear tear between 1 and 5")
	}

	asset := steam.NewAsset(assetName, wearTier, statTrak)

	marketListing, err := steam.GetMarketListing(asset.EncodedName)
	if err != nil {
		log.Fatalf("failed get market listing: %s", err)
	}

	// asset, err := steam.GetAsset()
	// if err != nil {
	// 	log.Fatal(err)
	// }

	marketListingJSON, err := json.MarshalIndent(marketListing, "", "\t")
	if err != nil {
		log.Fatalf("failed to marshal listing JSON: %s", err)
	}

	fmt.Printf("%s\n", marketListingJSON)

	// err := steam.ListAssets()
	// if err != nil {
	// 	log.Fatal(err)
	// }
}
