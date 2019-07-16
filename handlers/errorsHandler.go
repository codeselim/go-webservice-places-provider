package handlers

import (
	"app/api"
	"encoding/json"
	"log"
	"net/http"
)

// Customized error handler, checks types of returned error,
// and take the last action sending the right response accordingly
func HandleError(err error, w http.ResponseWriter, r *http.Request) {
	switch e := err.(type) {
	case *api.Error:
		w.WriteHeader(e.StatusCode)
		json.NewEncoder(w).Encode(e)

		//... extend here other cases and Errors types

	default:
		// Any error types we don't specifically look out for, defaults
		// to serving a HTTP 500 - Internal Server Error - writes a JSON response
		resp := api.Error{
			StatusCode: http.StatusInternalServerError,
			TraceId: GetRequestID(r.Context()),
			Message: "Oops! something went wrong! Please refer to our support with your traceId",
		}
		log.Printf("Error: %s  TraceId %s", e.Error(), GetRequestID(r.Context()))
		json.NewEncoder(w).Encode(resp)
	}
}
