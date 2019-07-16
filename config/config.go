package config

import (
	"log"
	"os"
	"time"
)

/**
 * Config package for a very basic configuration management.
 */
const (
	DefaultSearchRadius         = 100 //km
	DefaultGooglePlacesLanguage = "en"
	DefaultHttpServerPort       = "8081"
	DefaultProviderTimeout		= 10 * time.Second
)

var GooglePlacesApiKey = os.Getenv("GOOGLE_PLACES_API_KEY")
var FoursquareClientID = os.Getenv("FOURSQUARE_CLIENT_ID")
var FoursquareClientSecret = os.Getenv("FOURSQUARE_CLIENT_SECRET")

func init() {
	log.Print("configuration loaded...")
}
