package providers

import (
	"app/api"
	"app/config"
	"context"
	"fmt"
	"github.com/peppage/foursquarego"
	"log"
)

/**
 * Places provider: Foursquare
 * API ref: https://developer.foursquare.com/docs/api/venues/suggestcompletion
 * Provides places results from the Foursquare Venues API
 */

// since we are obliged ot use a location with the search. In case Lat/Lng are not supplied
// Just for this demo, in real life, once can adopt other solutions (e.g. location from IP...)
const (
	fallbackLat = 53.564615
	fallbackLng = 9.918173
)

type foursquareProvider struct {
	providerLabel ProviderLabel
	providerConfig ProviderConfig
}

// Constructor
func NewFoursquareProvider() Provider {
	return &foursquareProvider{
		providerLabel: FoursquareLabel,
	}
}

func (f *foursquareProvider) GetPlacesByQuery(ctx context.Context, request PlaceSearchRequest) (places api.Places, err error) {
	log.Print(f.providerConfig.Timeout)
	httpClient := getHttpClientFromConfig(&f.providerConfig)
	client := foursquarego.NewClient(httpClient, "foursquare", config.FoursquareClientID, config.FoursquareClientSecret, "")

	searchParam := &foursquarego.VenueSuggestParams{
		Radius: config.DefaultSearchRadius,
		Query:  request.InputString,
	}

	if request.Location != nil {
		searchParam.LatLong = fmt.Sprintf("%.6f,%.6f", request.Location.Lat, request.Location.Lng)
	} else {
		// since we are obliged ot use a location with the search. In case Lat/Lng are not supplied
		// Just for this demo, in real life, once can adopt other solutions (e.g. location from IP...)
		searchParam.LatLong = fmt.Sprintf("%.6f,%.6f", fallbackLat, fallbackLng)
	}
	// Get venues suggestions
	miniVenues, _, err := client.Venues.SuggestCompletion(searchParam)
	if err != nil {
		return api.Places{}, err
	}

	places = f.fourSquarePlacesToApiPlacesConverter(miniVenues)
	return places, nil
}

func (f *foursquareProvider) GetPlaceDetails(ctx context.Context, placeId string) (placeDetails api.PlaceDetails, err error) {
	//dummy implementation, extend in the future
	return api.PlaceDetails{ID: placeId, Name: "John Smith", SomeText: "Endpoint Not implemented yet! Take it easy!"}, nil
}

func (f *foursquareProvider) WithConfig(config ProviderConfig) Provider {
	f.providerConfig = config
	return f
}

// Converter Foursquare Models -> API Models
func (f *foursquareProvider) fourSquarePlacesToApiPlacesConverter(venues []foursquarego.MiniVenue) api.Places {
	places := api.Places{}
	for _, venue := range venues {
		place := api.Place{
			ID:       venue.ID,
			Name:     venue.Name,
			Provider: string(f.providerLabel),
			URI:      fmt.Sprintf("/fs/%s/details", venue.ID), //kind of hateoas href
			Address:  getFormattedAddress(venue),
			Location: &api.Location{
				Lat: venue.Location.Lat,
				Lng: venue.Location.Lng,
			},
		}
		places = append(places, place)
	}
	return places
}

func getFormattedAddress(venue foursquarego.MiniVenue) string {
	// Either a full address or nothing!
	// (since foursquare can return incomplete addresses sometimes )
	if venue.Location.Address != "" &&
		venue.Location.PostalCode != "" &&
		venue.Location.City != "" &&
		venue.Location.Country != "" {
		return fmt.Sprintf("%s, %s %s, %s", venue.Location.Address, venue.Location.PostalCode, venue.Location.City, venue.Location.Country)
	} else {
		return ""
	}
}

