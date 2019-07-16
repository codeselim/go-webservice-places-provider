package providers

import (
	"app/api"
	"app/config"
	"context"
	"fmt"
	"googlemaps.github.io/maps"
	"log"
)

/**
 * Places provider: Google Places
 * API ref: https://developers.google.com/places/web-service/autocomplete#place_autocomplete_results
 * Provides places results from the Google Places API
 */

type googlePlacesProvider struct {
	providerLabel ProviderLabel
	providerConfig ProviderConfig
}

// Constructor
func NewGoogleLocationProvider() Provider {
	return &googlePlacesProvider{
		providerLabel: GooglePlacesProviderLabel,
	}
}

func (g *googlePlacesProvider) GetPlacesByQuery(ctx context.Context, request PlaceSearchRequest) (places api.Places, err error) {
	log.Print(g.providerConfig.Timeout)

	var client *maps.Client
	httpClient := getHttpClientFromConfig(&g.providerConfig)
	apiKey := config.GooglePlacesApiKey
	client, err = maps.NewClient(maps.WithAPIKey(apiKey), maps.WithHTTPClient(httpClient))
	if err != nil {
		return api.Places{}, err
	}

	searchParam := &maps.PlaceAutocompleteRequest{
		Input:    request.InputString,
		Language: config.DefaultGooglePlacesLanguage,
		Radius:   uint(config.DefaultSearchRadius),
		Types:    "establishment",
	}

	if request.Location != nil {
		searchParam.Location = &maps.LatLng{
			Lat: request.Location.Lat,
			Lng: request.Location.Lng,
		}
	}

	resp, err := client.PlaceAutocomplete(context.Background(), searchParam)
	if err != nil {
		return api.Places{}, err
	}

	apiPlaces := g.googlePlacesToApiPlacesConverter(resp)
	return apiPlaces, nil
}

func (g *googlePlacesProvider) GetPlaceDetails(ctx context.Context, placeId string) (placeDetails api.PlaceDetails, err error) {
	//dummy implementation, extend in the future
	return api.PlaceDetails{ID: placeId, Name: "John Smith", SomeText: "Endpoint Not implemented yet! Take it easy!"}, nil
}

func (g *googlePlacesProvider) WithConfig(config ProviderConfig) Provider {
	g.providerConfig = config
	return g
}

// Converter Google Models -> API Models
func (g *googlePlacesProvider) googlePlacesToApiPlacesConverter(resp maps.AutocompleteResponse) api.Places {
	places := api.Places{}
	Predictions := resp.Predictions

	for _, prediction := range Predictions {
		place := api.Place{
			Provider: string(g.providerLabel),
			Address:  prediction.StructuredFormatting.SecondaryText, //relying on the secondary text since the search types is fixed on establishments
			Name:     prediction.StructuredFormatting.MainText,
			ID:       prediction.PlaceID,
			Location: nil,                                               // no place location details in the returned results
			URI:      fmt.Sprintf("/gp/%s/details", prediction.PlaceID), //kind of hateoas href
		}
		places = append(places, place)
	}
	return places
}
