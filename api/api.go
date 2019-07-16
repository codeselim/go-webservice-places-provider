package api

import (
	"fmt"
)

type Location struct {
	Lng float64 `json:"lng,omitempty"`
	Lat float64 `json:"lat,omitempty"`
}

type Place struct {
	ID       string    `json:"id"`
	Provider string    `json:"provider"`
	Name     string    `json:"name"`
	Location *Location `json:"location,omitempty"`
	Address  string    `json:"address,omitempty"`
	URI      string    `json:"uri"`
}

type Places []Place

type Error struct {
	StatusCode int    `json:"-"`                 // http status code. It will not be marshaled, instead used as a header
	TraceId    string `json:"traceId,omitempty"` // can be a tracing id/correlation id
	Type       string `json:"type,omitempty"`    // error type example "OAuthException" can be also an internal custom go Error type following the error interface
	Code       int    `json:"code,omitempty"`    // internal application code.
	Message    string `json:"message"`           // Human readable message
	// we can also introduce other fields like "retryable (bool)" for a client back off logic
}

// interface golang/error
func (e *Error) Error() string {
	return fmt.Sprintf("%#v", e)
}

const (
	TextInputParamIsMissingErrorCode = 10001
	LatLngParamMalformedErrorCode    = 10002
	//... can be extended in the future
)

var ErrorMessageText = map[int]string{
	TextInputParamIsMissingErrorCode: "the 'text' query parameter is missing",
	LatLngParamMalformedErrorCode:    "Malformed lat, lng parameters",
	//... can be extended in the future
}

/*
 * Just a holder for future implementations
 */
type PlaceDetails struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	SomeText string `json:"someText"`
	// All other details
}
