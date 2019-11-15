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
	defaultListings  int    = 25
)

var (
	assetName   string
	steamAPIKey string
	wearTier    int
	listings    int
	statTrak    bool
	debug       bool
)

func init() {
	flag.StringVar(&assetName, "n", defaultAssetName, "the name of the Steam asset to query")
	flag.StringVar(&steamAPIKey, "k", "", "the user Steam Web API Key")
	flag.IntVar(&wearTier, "w", defaultWearTier, "what wear quality to query")
	flag.IntVar(&listings, "l", defaultListings, "how many market listings, default 25")
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
		log.Fatalf("Failed to get asset listings %s", err)
	}

	if assetList == nil {
		log.Fatalf("No results for %s", assetName)
	}

	assetJSON, err := json.MarshalIndent(assetList, "", "\t")
	if err != nil {
		log.Fatalf("failed to marshal listing JSON: %s", err)
	}

	notableIDs := steam.CheckForRarity(*assetList)
	highlight := ""
	if len(notableIDs) > 0 {
		for _, id := range notableIDs {
			highlight += fmt.Sprintf("\nHIGHLIGHT: %s", id)
		}
	}

	fmt.Printf("%s\n\n%s", assetJSON, highlight)
}
