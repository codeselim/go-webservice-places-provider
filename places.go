package main

import (
	"flag"
	"github.com/codeselim/go-webservice-places-provider/config"
	"github.com/codeselim/go-webservice-places-provider/handlers"
	"github.com/codeselim/go-webservice-places-provider/log"
	"github.com/codeselim/go-webservice-places-provider/providers"
	gh "github.com/gorilla/handlers"
	"github.com/gorilla/mux"

	"net/http"
	"time"
)

const apiVersion = "v1"

var webServerPort string

func init() {
	flag.StringVar(&webServerPort, "httpServerPort", config.DefaultHttpServerPort, "Default port to expose on the API. use -httpServerPort=<port_value>")
}

func main() {
	// parse app flags
	flag.Parse()
	logger := log.GetLogger()

	// Bootstrap the application
	// Providers
	googlePlacesConfig := providers.ProviderConfig{Timeout: time.Second * 12, Language: "en"} //else config will fall to defaults
	foursquareConfig := providers.ProviderConfig{Timeout: time.Second * 13}                   // for example...
	googlePlacesProvider := providers.NewGoogleLocationProvider(&googlePlacesConfig)
	foursquareProvider := providers.NewFoursquareProvider(&foursquareConfig)
	placesHandler := handlers.NewPlacesHandler(googlePlacesProvider, foursquareProvider) //extend and provide as many providers as you want!

	// Other handlers
	recoveryHandler := gh.RecoveryHandler()
	// todo compress handler ...etc

	r := mux.NewRouter()
	r.Use(handlers.RequestIdMiddleware, handlers.LoggingMiddleware)
	r.HandleFunc("/api/"+apiVersion+"/places", placesHandler.GetPlaces).Methods("GET")
	r.HandleFunc("/api/"+apiVersion+"/status", handlers.GetStatus).Methods("GET")
	logger.Info("Serving requests on port: " + webServerPort)
	logger.Fatal(http.ListenAndServe(":"+webServerPort, recoveryHandler(r)))
}
