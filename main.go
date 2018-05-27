package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"

	"./steam"
)

const (
	defaultAssetName = "Chroma 3 Case"
)

var (
	assetName   string
	steamAPIKey string
)

func init() {
	flag.StringVar(&assetName, "name", defaultAssetName, "the name of the Steam asset to lookup")
	flag.StringVar(&steamAPIKey, "key", "", "the user Steam Web API Key")
	flag.Parse()
}

func main() {
	if steamAPIKey == "" {
		log.Fatal("please specify API Key")
	}

	asset := steam.NewAsset(assetName)

	if asset.EncodedName == "" {
		log.Fatalf("encoded name blank: %s", asset.EncodedName)
	}

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
