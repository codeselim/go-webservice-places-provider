package providers

import (
	"github.com/peppage/foursquarego"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestUnitNewFoursquareProvider(t *testing.T) {
	assert.NotPanics(t, func() { NewFoursquareProvider(&ProviderConfig{}) })
}

func TestUnitFSWithConfig(t *testing.T) {
	//empty config
	assert.NotPanics(t, func() { NewFoursquareProvider(&ProviderConfig{}) })
	//filled config
	assert.NotPanics(t, func() { NewFoursquareProvider(&ProviderConfig{Timeout: time.Second * 10}) })
}

func TestUnitfourSquarePlacesToApiPlacesConverter(t *testing.T) {
	venues := []foursquarego.MiniVenue{
		{
			Name: "name here",
			Location: foursquarego.Location{
				Lng: 13.3,
				Lat: 32.32,
			},
			ID: "ID-here",
			//...
		},
		{
			Name: "another name here",
			Location: foursquarego.Location{
				Lng: 11.3,
				Lat: 32.32,
			},
			ID: "another-ID",
			//...
		},
	}

	actualApiPlaces := fourSquarePlacesToApiPlacesConverter(venues)

	//additional structure verifications can be added... e.g. with reflect.DeepEqual
	assert.Equal(t, 2, len(actualApiPlaces))
}

func TestUnitgetFormattedAddress(t *testing.T) {
	venue := foursquarego.MiniVenue{
		Name: "name here",
		Location: foursquarego.Location{
			Lng: 13.3,
			Lat: 32.32,
		},
		ID: "ID-here",
		//...
	}

	address := getFormattedAddress(venue)
	assert.Equal(t, "", address)

	venue.Location.Country = "DE"
	venue.Location.City = "Hamburg"

	address = getFormattedAddress(venue)
	assert.Equal(t, "", address)
}
