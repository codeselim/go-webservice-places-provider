package handlers

import (
	"encoding/json"
	"github.com/codeselim/go-webservice-places-provider/api"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

// This is the demo how to test a handler function and assert on the json body
func TestUnitGetStatus(t *testing.T) {
	// create request
	req, err := http.NewRequest("GET", "/api/v1/status", nil)
	if err != nil {
		t.Fatal(err)
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetStatus)

	// The handler satisfy http.Handler, so we can call its ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned unexpected status code: got %v want %v",
			status, http.StatusBadRequest)
	}

	expectedStatus := api.Status{
		Message:    apiStatusMessage,
		StateLabel: apiStatusStateLabel,
	}
	expectedStatusJson, err := json.Marshal(expectedStatus)
	if err != nil {
		t.Fatal(err)
	}
	assert.JSONEq(t, string(expectedStatusJson), rr.Body.String())
}
