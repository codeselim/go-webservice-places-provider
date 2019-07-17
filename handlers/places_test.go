package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/codeselim/go-webservice-places-provider/api"
	"github.com/codeselim/go-webservice-places-provider/providers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"testing"
)

var (
	apiPlaceFromGoogle = api.Place{
		ID:       "id1",
		Address:  "address1",
		Name:     "place1",
		URI:      "Uri1",
		Provider: "google-provider-label",
	}

	apiPlaceFromFoursquare = api.Place{
		ID:       "id2",
		Address:  "address2",
		Name:     "place2",
		URI:      "Uri2",
		Provider: "foursquare-provider-label",
		Location: &api.Location{
			Lat: 32.22,
			Lng: 16.77,
		},
	}
)

type mockGooglePlacesProvider struct {
	mock.Mock
}

//Meet the Provider interface
func (m *mockGooglePlacesProvider) GetPlacesByQuery(ctx context.Context, request providers.PlaceSearchRequest) (api.Places, error) {
	args := m.Called(ctx, request)
	return args.Get(0).(api.Places), args.Error(1)
}
func (m *mockGooglePlacesProvider) GetPlaceDetails(ctx context.Context, placeId string) (api.PlaceDetails, error) {
	args := m.Called(ctx, placeId)
	return args.Get(0).(api.PlaceDetails), args.Error(1)
}
func (m *mockGooglePlacesProvider) GetProviderLabel() providers.ProviderLabel {
	return "google-provider-label"
}

type mockFourSquarePlacesProvider struct {
	mock.Mock
}

//Meet the Provider interface
func (m *mockFourSquarePlacesProvider) GetPlacesByQuery(ctx context.Context, request providers.PlaceSearchRequest) (api.Places, error) {
	args := m.Called(ctx, request)
	return args.Get(0).(api.Places), args.Error(1)
}
func (m *mockFourSquarePlacesProvider) GetPlaceDetails(ctx context.Context, placeId string) (api.PlaceDetails, error) {
	args := m.Called(ctx, placeId)
	return args.Get(0).(api.PlaceDetails), args.Error(1)
}
func (m *mockFourSquarePlacesProvider) GetProviderLabel() providers.ProviderLabel {
	return "foursquare-provider-label"
}

func TestPlacesHandlerPanics(t *testing.T) {
	assert.Panics(t, func() { NewPlacesHandler(nil) })
	mockedGooglePlacesProvider := new(mockGooglePlacesProvider)
	assert.NotPanics(t, func() { NewPlacesHandler(mockedGooglePlacesProvider) })
}

func TestPlacesHandlerGetPlacesEmptyResponse(t *testing.T) {
	// Create a request to pass to our handler.
	req, err := http.NewRequest("GET", "/api/v1/places?text=chocolate", nil)
	if err != nil {
		t.Fatal(err)
	}

	googlePlacesProvider := new(mockGooglePlacesProvider)
	googlePlacesProvider.On("GetPlacesByQuery", mock.Anything, mock.Anything).Return(api.Places{}, nil)
	placesHandler := NewPlacesHandler(googlePlacesProvider)

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(placesHandler.GetPlaces)

	// The handler satisfy http.Handler, so we can call its ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned unexpected status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	expected := `[]`
	if rr.Body.String() != expected {
		assert.JSONEq(t, expected, rr.Body.String())
	}
}

func TestPlacesHandlerGetPlacesErrorMissingParams(t *testing.T) {
	// Create a request with missing parameters.
	req, err := http.NewRequest("GET", "/api/v1/places", nil)
	if err != nil {
		t.Fatal(err)
	}

	googlePlacesProvider := new(mockGooglePlacesProvider)
	googlePlacesProvider.On("GetPlacesByQuery", mock.Anything, mock.Anything).Return(api.Places{apiPlaceFromGoogle}, nil)
	placesHandler := NewPlacesHandler(googlePlacesProvider)

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(placesHandler.GetPlaces)

	// The handler satisfy http.Handler, so we can call its ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned unexpected status code: got %v want %v",
			status, http.StatusBadRequest)
	}

}

func TestPlacesHandlerGetPlacesFromAllProviders(t *testing.T) {
	// create request
	req, err := http.NewRequest("GET", "/api/v1/places?text=chocolate", nil)
	if err != nil {
		t.Fatal(err)
	}

	googlePlacesProvider := new(mockGooglePlacesProvider)
	foursquarePlacesProvider := new(mockFourSquarePlacesProvider)

	googlePlacesProvider.On("GetPlacesByQuery", mock.Anything, mock.Anything).Return(api.Places{apiPlaceFromGoogle}, nil)
	foursquarePlacesProvider.On("GetPlacesByQuery", mock.Anything, mock.Anything).Return(api.Places{apiPlaceFromFoursquare}, nil)

	placesHandler := NewPlacesHandler(googlePlacesProvider, foursquarePlacesProvider)

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(placesHandler.GetPlaces)

	// The handler satisfy http.Handler, so we can call its ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned unexpected status code: got %v want %v",
			status, http.StatusBadRequest)
	}

	//Unmarshal json body in structs and assert!
	actualPlaces := api.Places{}

	err = json.Unmarshal(rr.Body.Bytes(), &actualPlaces)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, 2, len(actualPlaces))

	for _, actualPlace := range actualPlaces {
		if actualPlace.Provider == "google-provider-label" {
			assert.Equal(t, apiPlaceFromGoogle, actualPlace)
		}
		if actualPlace.Provider == "foursquare-provider-label" {
			assert.Equal(t, apiPlaceFromFoursquare, actualPlace)
		}
	}
}

func TestPlacesHandlerGetPlacesFromAllProvidersSkipOneProviderErr(t *testing.T) {
	// create request
	req, err := http.NewRequest("GET", "/api/v1/places?text=chocolate", nil)
	if err != nil {
		t.Fatal(err)
	}

	googlePlacesProvider := new(mockGooglePlacesProvider)
	foursquarePlacesProvider := new(mockFourSquarePlacesProvider)

	errorInProvider := errors.New("some-kind-of-error")

	googlePlacesProvider.On("GetPlacesByQuery", mock.Anything, mock.Anything).Return(api.Places{}, errorInProvider)
	foursquarePlacesProvider.On("GetPlacesByQuery", mock.Anything, mock.Anything).Return(api.Places{apiPlaceFromFoursquare}, nil)

	placesHandler := NewPlacesHandler(googlePlacesProvider, foursquarePlacesProvider)

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(placesHandler.GetPlaces)

	// The handler satisfy http.Handler, so we can call its ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned unexpected status code: got %v want %v",
			status, http.StatusBadRequest)
	}

	//Unmarshal json body in structs and assert!
	actualPlaces := api.Places{}

	err = json.Unmarshal(rr.Body.Bytes(), &actualPlaces)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, 1, len(actualPlaces))

	assert.Equal(t, api.Places{apiPlaceFromFoursquare}, actualPlaces)
}
