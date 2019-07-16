package main

import (
	"app/config"
	"app/handlers"
	"app/providers"
	"flag"
	gh "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"log"
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

	// Bootstrap the application
	// Providers
	googlePlacesConfig := providers.ProviderConfig{ Timeout: time.Second * 12 } //else config will fall to defaults
	foursquareConfig := providers.ProviderConfig{ Timeout: time.Second * 13 } // for example...
	googlePlacesProvider := providers.NewGoogleLocationProvider().WithConfig(googlePlacesConfig)
	foursquareProvider := providers.NewFoursquareProvider().WithConfig(foursquareConfig)
	placesHandler := handlers.NewPlacesHandler(googlePlacesProvider, foursquareProvider) //extend and provide as many providers as you want!

	// Other handlers
	recoveryHandler := gh.RecoveryHandler()
	// todo compress handler ...etc

	r := mux.NewRouter()
	r.Use(handlers.RequestIdMiddleware, handlers.LoggingMiddleware)
	r.HandleFunc("/"+apiVersion+"/places", placesHandler.GetPlaces).Methods("GET")
	r.HandleFunc("/"+apiVersion+"/status", placesHandler.GetStatus).Methods("GET")
	log.Print("Serving requests on port: " + webServerPort)
	log.Fatal(http.ListenAndServe(":"+webServerPort, recoveryHandler(r)))
}
