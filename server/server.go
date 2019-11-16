package server

import (
	"encoding/json"
	"flag"
	"html/template"
	"log"
	"net/http"

	"github.com/jtrotsky/eiffel65/steam"
)

const (
	defaultAssetName        string = "AK-47 | Case Hardened"
	defaultWearTier         int    = 3 // Field-Tested
	defaultNumberOfListings int    = 25
)

// Config contains configuration details and search arguments
type Config struct {
	SteamAPIKey      string
	NumberOfListings int
	Debug            bool
	AssetName        string
	WearTier         int
	StatTrak         bool
}

var config = &Config{
	SteamAPIKey:      "",
	NumberOfListings: defaultNumberOfListings,
	AssetName:        defaultAssetName,
	WearTier:         defaultWearTier,
	StatTrak:         false,
}

func getConfig(config *Config) {
	flag.StringVar(&config.AssetName, "n", defaultAssetName, "the name of the Steam asset to query")
	flag.StringVar(&config.SteamAPIKey, "k", "", "the user Steam Web API Key")
	flag.IntVar(&config.WearTier, "w", defaultWearTier, "what wear quality to query")
	flag.IntVar(&config.NumberOfListings, "l", defaultNumberOfListings, "how many market listings, default 25")
	flag.BoolVar(&config.StatTrak, "s", false, "whether to query items with StatTrak")
	flag.BoolVar(&config.Debug, "d", false, "debug mode")
	flag.Parse()
}

// Start hosts a server on localhost:8080
func Start() error {
	getConfig(config)

	if config.SteamAPIKey == "" {
		log.Fatal("please specify an API Key")
	}

	if config.WearTier > 5 || config.WearTier < 1 {
		log.Fatal("please specify a wear tear between 1 and 5")
	}

	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/favicon.ico", faviconHandler)
	http.HandleFunc("/listings", listingsHandler)

	return http.ListenAndServe(":8080", nil)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	templates, err := template.ParseFiles("public/index.html")
	if err != nil {
		log.Println("failed to parse html templates")
	}

	templates.ExecuteTemplate(w, "index.html", nil)
}

func faviconHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "public/favicon.ico")
}

func listingsHandler(w http.ResponseWriter, r *http.Request) {
	steamClient := steam.NewClient(config.SteamAPIKey)

	assetList, err := steamClient.NewAsset(config.AssetName, config.WearTier,
		config.NumberOfListings, config.StatTrak, config.Debug)
	if err != nil {
		log.Fatalf("Failed to get asset listings %s", err)
	}

	if assetList == nil {
		log.Fatalf("No results for %s", config.AssetName)
	}

	assetJSON, err := json.MarshalIndent(assetList, "", "\t")
	if err != nil {
		log.Fatalf("failed to marshal listing JSON: %s", err)
	}

	// notableIDs := steam.CheckForRarity(*assetList)
	// highlight := ""
	// if len(notableIDs) > 0 {
	// 	for id, asset := range notableIDs {
	// 		highlight += fmt.Sprintf("\nHIGHLIGHT: %s SEED: %d PRICE: %s%s SCREENSHOT: %s",
	// 			id, asset.Float.PaintSeed, asset.ListingTotalPrice, asset.ListingCurrency, asset.ScreenshotURL)
	// 	}
	// }

	w.Write(assetJSON)
}
