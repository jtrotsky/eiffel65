package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/jtrotsky/eiffel65/steam"
)

const (
	// defaultAssetName string = "AK-47 | Case Hardened"
	defaultAssetName string = "Case Hardened"
	defaultWearTier  int    = 3
	defaultListings  int    = 25
)

var (
	assetName   string
	steamAPIKey string
	wearTier    int
	listings    int
	statTrak    bool
	localbrowse bool
	debug       bool
)

func init() {
	flag.StringVar(&assetName, "n", defaultAssetName, "the name of the Steam asset to query")
	flag.StringVar(&steamAPIKey, "k", "", "Steam Web API Key")
	flag.IntVar(&wearTier, "w", defaultWearTier, "what weapon wear quality: 1 is Factory New, 5 is Battle-Scarred, the default 3 Field-Tested")
	flag.IntVar(&listings, "l", defaultListings, "how many market listings, default 25")
	flag.BoolVar(&statTrak, "s", false, "whether to query items with StatTrak")
	flag.BoolVar(&localbrowse, "x", false, "browser mode")
	flag.BoolVar(&debug, "d", false, "debug mode")
	flag.Parse()
}

func main() {
	if steamAPIKey == "" {
		log.Fatal("please specify an API Key with -k")
	}

	switch {
	case wearTier < 1:
		wearTier = 1
	case wearTier > 5:
		wearTier = 5
	default:
	}

	steamClient := steam.NewClient(steamAPIKey, debug)

	assetList, err := steamClient.NewAsset(assetName, wearTier, listings, statTrak)
	if err != nil {
		log.Fatal(err)
	}

	if assetList == nil {
		log.Fatal(err)
	}

	assetJSON, err := json.MarshalIndent(assetList, "", "\t")
	if err != nil {
		log.Fatal(err)
	}

	notableIDs := steam.CheckForRarity(*assetList)
	highlight := ""
	if len(notableIDs) > 0 {
		for id, asset := range notableIDs {
			highlight += fmt.Sprintf("\nHIGHLIGHT: %s SEED: %d FLOAT: %.4f PRICE: %s%s SCREENSHOT: %s",
				id, asset.Float.PaintSeed, asset.Float.FloatValue, asset.ListingTotalPrice, asset.ListingCurrency, asset.ScreenshotURL)
		}
	}

	if localbrowse {
		http.HandleFunc("/", handler)
		log.Fatal(http.ListenAndServe(":8080", nil))
	} else {
		fmt.Printf("%s\n\n%s", assetJSON, highlight)
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome %s", r.URL.Path[1:])
}
