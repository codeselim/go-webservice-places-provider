package handlers

import (
	"encoding/json"
	"github.com/codeselim/go-webservice-places-provider/api"
	"github.com/codeselim/go-webservice-places-provider/log"

	"net/http"
)

// Customized error handler, checks types of returned error,
// and take the last action sending the right response accordingly
func HandleError(err error, w http.ResponseWriter, r *http.Request) {

	loggerWithContext := log.GetLoggerWithContext(r.Context())

	switch e := err.(type) {
	case *api.Error:
		w.WriteHeader(e.StatusCode)
		json.NewEncoder(w).Encode(e)
		loggerWithContext.Error(e.Error())

		//... extend here other cases and Errors types

	default:
		// Any error types we don't specifically look out for, defaults
		// to serving a HTTP 500 - Internal Server Error - writes a JSON response
		resp := api.Error{
			StatusCode: http.StatusInternalServerError,
			TraceId:    GetRequestID(r.Context()),
			Message:    "Oops! something went wrong! Please refer to our support with your traceId",
		}
		loggerWithContext.Error(e.Error())
		json.NewEncoder(w).Encode(resp)
	}
}
