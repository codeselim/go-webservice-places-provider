package handlers

import (
	"encoding/json"
	"github.com/codeselim/go-webservice-places-provider/api"
	"net/http"
)

const (
	apiStatusMessage    = "API v1 Alive!"
	apiStatusStateLabel = "READY"
)

func GetStatus(w http.ResponseWriter, req *http.Request) {
	status := api.Status{
		Message:    apiStatusMessage,
		StateLabel: apiStatusStateLabel,
	}
	setDefaultHeaders(w) // should be done before writeHeader
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(status)
}
