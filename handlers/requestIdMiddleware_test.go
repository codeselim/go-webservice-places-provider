package handlers

import (
	"github.com/gorilla/mux"
	_ "github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func dummyHandler(w http.ResponseWriter, r *http.Request) {}

func TestUnitRequestIdMiddleware(t *testing.T) {

	router := mux.NewRouter()

	rw := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	router.Use(RequestIdMiddleware)
	router.HandleFunc("/", dummyHandler).Methods("GET")
	router.ServeHTTP(rw, req)

	//reqId := GetRequestID(req.Context())

	//assert.NotEmpty(t, reqId)
}
