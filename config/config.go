package config

import (
	"log"
	"os"
	"sync"
	"time"
)

/**
 * Config package for a very basic configuration management.
 */

const (
	DefaultSearchRadius         = 100 //km
	DefaultGooglePlacesLanguage = "en"
	DefaultHttpServerPort       = "8081"
	DefaultProviderTimeout      = 10 * time.Second
	MaxAllowedSearchRadius      = 50
	DefaultLoggingLevel         = "info"
)

type configSchema struct {
	GooglePlacesApiKey     string
	FoursquareClientID     string
	FoursquareClientSecret string
	DefaultHttpHeaders     map[string]string
	//sync.RWMutex : Mutexes can be added if config would be extended to add write actions
}

var (
	c    *configSchema
	once sync.Once
)

func Config() *configSchema {
	once.Do(func() { //https://golang.org/src/sync/once.go?s=1137:1164#L25
		c = &configSchema{
			GooglePlacesApiKey:     os.Getenv("GOOGLE_PLACES_API_KEY"),
			FoursquareClientID:     os.Getenv("FOURSQUARE_CLIENT_ID"),
			FoursquareClientSecret: os.Getenv("FOURSQUARE_CLIENT_SECRET"),
			DefaultHttpHeaders: map[string]string{
				"Access-Control-Allow-Origin": "*",
			},
		}
	})
	return c
}

func init() {
	log.Print("Configuration loaded...")
}
