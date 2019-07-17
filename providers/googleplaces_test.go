package providers

import (
	"github.com/stretchr/testify/assert"
	"googlemaps.github.io/maps"
	"testing"
	"time"
)

func TestUnitNewGoogleLocationProvider(t *testing.T) {
	assert.NotPanics(t, func() { NewGoogleLocationProvider(&ProviderConfig{}) })
	//check label
}

func TestUnitGPWithConfig(t *testing.T) {
	//empty config
	assert.NotPanics(t, func() { NewGoogleLocationProvider(&ProviderConfig{}) })
	//filled config
	assert.NotPanics(t, func() { NewGoogleLocationProvider(&ProviderConfig{Timeout: time.Second * 10}) })
}

func TestUnitgooglePlacesToApiPlacesConverter(t *testing.T) {
	googlePredictions := []maps.AutocompletePrediction{
		{
			PlaceID: "someID",
			StructuredFormatting: maps.AutocompleteStructuredFormatting{
				MainText:      "Main text here",
				SecondaryText: "Second text here",
			},
		},
		{
			PlaceID: "someID2",
			StructuredFormatting: maps.AutocompleteStructuredFormatting{
				MainText:      "More text here",
				SecondaryText: "even more text here",
			},
		},
	}

	googleResp := maps.AutocompleteResponse{
		Predictions: googlePredictions,
	}

	actualApiPlaces := googlePlacesToApiPlacesConverter(googleResp)

	//additional structure verifications can be added... e.g. with reflect.DeepEqual
	assert.Equal(t, 2, len(actualApiPlaces))
}
