package providers

import (
	"context"
	"fmt"
	"github.com/codeselim/go-webservice-places-provider/api"
	"github.com/codeselim/go-webservice-places-provider/config"
	"github.com/codeselim/go-webservice-places-provider/log"
	"googlemaps.github.io/maps"
)

/**
 * Places provider: Google Places
 * API ref: https://developers.google.com/places/web-service/autocomplete#place_autocomplete_results
 * Provides places results from the Google Places API
 */

type googlePlacesProvider struct {
	providerLabel  ProviderLabel
	providerConfig *ProviderConfig
	mapsClient     *maps.Client
}

// Constructor
func NewGoogleLocationProvider(providerConfig *ProviderConfig) Provider {
	if providerConfig == nil {
		log.GetLogger().Panic("ProviderConfig should be provided")
	}

	httpClient := getHttpClientFromConfig(providerConfig)
	apiKey := config.Config().GooglePlacesApiKey
	client, err := maps.NewClient(maps.WithAPIKey(apiKey), maps.WithHTTPClient(httpClient))

	if err != nil {
		log.GetLogger().Panic("Couldn't create maps client, are you sure you provided the needed Google credentials?")
	}

	return &googlePlacesProvider{
		providerLabel:  GooglePlacesProviderLabel,
		providerConfig: providerConfig,
		mapsClient:     client,
	}
}

func (g *googlePlacesProvider) GetPlacesByQuery(ctx context.Context, request PlaceSearchRequest) (places api.Places, err error) {
	language := config.DefaultGooglePlacesLanguage
	if g.providerConfig.Language != "" {
		language = g.providerConfig.Language
	}

	searchParam := &maps.PlaceAutocompleteRequest{
		Input:    request.InputString,
		Language: language,
		Radius:   uint(getSearchRadiusFromConfig(g.providerConfig)),
		Types:    "establishment",
	}

	if request.Location != nil {
		searchParam.Location = &maps.LatLng{
			Lat: request.Location.Lat,
			Lng: request.Location.Lng,
		}
	}

	resp, err := g.mapsClient.PlaceAutocomplete(context.Background(), searchParam)
	if err != nil {
		return api.Places{}, err
	}

	apiPlaces := googlePlacesToApiPlacesConverter(resp)
	return apiPlaces, nil
}

func (g *googlePlacesProvider) GetPlaceDetails(ctx context.Context, placeId string) (placeDetails api.PlaceDetails, err error) {
	//dummy implementation, extend in the future
	return api.PlaceDetails{ID: placeId, Name: "John Smith", SomeText: "Endpoint Not implemented yet! Take it easy!"}, nil
}

func (g *googlePlacesProvider) GetProviderLabel() ProviderLabel {
	return g.providerLabel
}

// Converter Google Models -> API Models
func googlePlacesToApiPlacesConverter(resp maps.AutocompleteResponse) api.Places {
	places := api.Places{}
	Predictions := resp.Predictions

	for _, prediction := range Predictions {
		place := api.Place{
			Provider: string(GooglePlacesProviderLabel),
			Address:  prediction.StructuredFormatting.SecondaryText, //relying on the secondary text since the search types is fixed on establishments
			Name:     prediction.StructuredFormatting.MainText,
			ID:       prediction.PlaceID,
			Location: nil,                                               // no place location details in the returned results
			URI:      fmt.Sprintf("/gp/details/%s", prediction.PlaceID), //kind of hateoas href
		}
		places = append(places, place)
	}
	return places
}
