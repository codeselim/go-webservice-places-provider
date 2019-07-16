package providers

import (
	"app/api"
	"app/config"
	"context"
	"net/http"
	"time"
)

type ProviderLabel string

const (
	GooglePlacesProviderLabel = ProviderLabel("GOOGLE_PLACES")
	FoursquareLabel           = ProviderLabel("FOURSQUARE")
	// ...
)

type ProviderConfig struct {
	Timeout time.Duration //configures a timeout to short-circuits long-running connections
	//... extend following requirements
}

type PlaceSearchRequest struct {
	InputString string
	Location    *Location
}

type Location struct {
	Lat float64
	Lng float64
}

type Provider interface {
	GetPlacesByQuery(ctx context.Context, request PlaceSearchRequest) (api.Places, error)
	GetPlaceDetails(ctx context.Context, placeId string) (api.PlaceDetails, error)
	WithConfig(config ProviderConfig) Provider
	//... extend interface
}

//helper function
func getHttpClientFromConfig(providerConfig *ProviderConfig) *http.Client {
	timeout := config.DefaultProviderTimeout
	if providerConfig != nil {
		timeout = providerConfig.Timeout
	}
	return &http.Client{
		Timeout: timeout,
	}
}
