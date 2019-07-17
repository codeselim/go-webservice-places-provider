package providers

import (
	"context"
	"github.com/codeselim/go-webservice-places-provider/api"
	"github.com/codeselim/go-webservice-places-provider/config"
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
	Timeout      time.Duration //configures a timeout to short-circuits long-running connections
	Language     string
	SearchRadius int //can be also provided as query param instead of internal config
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
	GetProviderLabel() ProviderLabel
	//... extend interface
}

//helper functions
func getHttpClientFromConfig(providerConfig *ProviderConfig) *http.Client {
	timeout := config.DefaultProviderTimeout
	if providerConfig != nil && providerConfig.Timeout > 0 {
		timeout = providerConfig.Timeout
	}
	return &http.Client{
		Timeout: timeout,
	}
}

func getSearchRadiusFromConfig(providerConfig *ProviderConfig) int {
	radius := config.DefaultSearchRadius
	if providerConfig.SearchRadius != 0 && providerConfig.SearchRadius < config.MaxAllowedSearchRadius {
		radius = providerConfig.SearchRadius
	}
	return radius
}
