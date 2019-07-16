package handlers

import (
	"app/api"
	"app/providers"
	"context"
	"encoding/json"
	"fmt"
	"golang.org/x/sync/errgroup"
	"net/http"
	"strconv"
)

type PlacesHandler struct {
	placesProviders []providers.Provider
}

func NewPlacesHandler(providers... providers.Provider) PlacesHandler{
	for _, provider := range providers {
		if provider == nil {
			panic("supplied providers cannot be nil")
		}
	}
	return PlacesHandler{
		placesProviders: providers,
	}
}

func (p *PlacesHandler) GetPlaces(w http.ResponseWriter, r *http.Request) {
	keys := r.URL.Query()
	inputString := keys.Get("text")
	lat := keys.Get("latitude")
	lng := keys.Get("longitude")

	if inputString == "" {
		apiError := &api.Error{
			Code:       api.TextInputParamIsMissingErrorCode,
			Message:    api.ErrorMessageText[api.TextInputParamIsMissingErrorCode],
			StatusCode: http.StatusBadRequest,
			TraceId:    GetRequestID(r.Context()),
		}
		HandleError(apiError, w, r)
		return
	}

	providerRequest := providers.PlaceSearchRequest{
		InputString: inputString,
	}

	if lat != "" && lng != "" {
		lat, errLat := strconv.ParseFloat(lat, 64)
		lng, errLng := strconv.ParseFloat(lng, 64)
		if errLat != nil || errLng != nil {
			apiError := &api.Error{
				Code:       api.LatLngParamMalformedErrorCode,
				Message:    api.ErrorMessageText[api.LatLngParamMalformedErrorCode],
				StatusCode: http.StatusBadRequest,
				TraceId:    GetRequestID(r.Context()),
			}
			HandleError(apiError, w, r)
			return
		}

		providerRequest.Location = &providers.Location{
			Lat: lat,
			Lng: lng,
		}
	}

	places, err := p.getPlacesParallel(r.Context(), providerRequest)
	if err != nil {
		HandleError(err, w, r)
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(places)
}

func (p *PlacesHandler) GetStatus(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "ALIVE!")
}

// Parallel execution of providers queries with sync/errGroup for a better error handling
// Ref. https://godoc.org/golang.org/x/sync/errgroup#ex-Group--Parallel
func (p *PlacesHandler) getPlacesParallel(ctx context.Context, request providers.PlaceSearchRequest) (api.Places, error) {
	g, ctx := errgroup.WithContext(ctx)
	placesResults := api.Places{}

	for _, provider := range p.placesProviders {
		provider := provider //check https://golang.org/doc/faq#closures_and_goroutines
		g.Go(func() error {
			places, err := provider.GetPlacesByQuery(ctx, request)
			if err == nil {
				placesResults = append(placesResults, places...)
			}
			return err
		})
	}
	if err := g.Wait(); err != nil {
		return nil, err
	}
	return placesResults, nil
}
