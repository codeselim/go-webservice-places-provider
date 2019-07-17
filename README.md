# Coding challenge - Webservice middleware application


This service, written in [Golang](https://golang.org/), exposes a RESTful API, which provides Places search functionality matching a consumer search parameters. 

This service is written using only basic Go libraries. No heavy frameworks are used, in order to demonstrate different part and concepts of a middlware webservice. 

This service is a middleware applications that fetches places from different providers:

* Google places's [autocomplete endpoint](https://developers.google.com/places/web-service/autocomplete#place_autocomplete_results) 
* Foursquare's [venues search endpoint](https://developer.foursquare.com/docs/api/venues/search) 

The service is extendable and other providers can be added easily. Refer to the Application Internals section. 
 
## Building the service

You need docker installed. The code is shipped with a Docker file responsible for a multi-stage container build. It will build an container containing the application binary statically compiled with all need libraries built in.   

Build an application image locally by running: 

`cd app`

`docker build . -t places-service` 

## Running the service

You can run the service containerized after performing the build described in the previous step. 

The service requires the presence of the following environment variables:

* GOOGLE_PLACES_API_KEY : should contain the Google places API Key 
* FOURSQUARE_CLIENT_ID : should contain the Foursquare client ID 
* FOURSQUARE_CLIENT_SECRET : should contain the Foursquare client secret 

Using vanilla Docker you can run it as follows (after you have built your image):

`docker run  -p 8081:8081 -e GOOGLE_PLACES_API_KEY='...' -e FOURSQUARE_CLIENT_ID='...' -e FOURSQUARE_CLIENT_SECRET='...' places-service`

Make sure that you fill your keys' values correctly. Look into the container's logs (stdout) for more info, in case of a malfunction.

By the fault, The container exposes the service on port **8081**. However you can override this behavior by exposing your favorite port and supplying the application the -httpServerPort flag

`docker run  -p <host_external_port>:<internal_port> --expose <internal_port> -e GOOGLE_PLACES_API_KEY='...' -e FOURSQUARE_CLIENT_ID='...' -e FOURSQUARE_CLIENT_SECRET='...' places-service -httpServerPort <internal_port>`      

## API usage

status endpoint
-------------- 
Request **GET host:port/v1/status**

Returns the API status. Such endpoints can be used for readiness/liveness probes. 

Answer:

Status code 200

```$xslt
{
"message": (string),
"state" (string),  
}
```
 
"state" can be only have the value "ACTIVE" for the time being. 

places endpoint
--------------
Request **GET host:port/v1/places**

### Query Parameters

| parameter | example | Description |
| :---: | :---:   | :---:       |
| text  | vegan   | **required** a search term to be applied against places names |
| latitude  | 53.6207518   | **optional** latitude of the user’s location (should be combined with longitude parameter).  |
| longitude  | 9.9881764   | **optional** longitude of the user’s location (should be combined with latitude parameter).   |

### Responses 
Content-Type : application/json

Possible responses:

*  *Success* : Status code 200 : Array(Place)
*  *Bad Request* : Status code 400 : Error
*  *Internal Server Error* : Status code 500 : Error

Place: 
```
{
    id:         (string) place id,
    provider:   (string) provider's label,
    name:       (string) place's name,
    address:    (string) place address - if applicable,
    location: {
        lng:    (number) the place longitude,
        lat:    (number) the place latitude
    },
    uri:        (string) URI of the place where more details are available
}
```

Error:
```
{
	traceId: (string)  Can be a tracing id/correlation id
	type:    (string)  Error type example "OAuthException"
	code:    (int)     Internal application code.
	message: (string)  Human readable message
}	
```

### Example call

text : car rental

latitude: 53.6207518

longitude: 9.9881764

GET http://localhost:8081/v1/places?text=car%20rental&latitude=53.6207518&longitude=9.9881764

Response: 

```json
[
...
{
  "id": "ChIJyxBFxlyIsUcRyuQUT68JAyI",
  "provider": "GOOGLE_PLACES",
  "name": "Enterprise Rent-A-Car - Flughafen Hamburg",
  "address": "Flughafenstraße, Hamburg, Germany",
  "uri": "/gp/ChIJyxBFxlyIsUcRyuQUT68JAyI/details"
},
{
  "id": "5111fd35e4b0752e2e6d7219",
  "provider": "FOURSQUARE",
  "name": "Rental Car Return",
  "location": {
    "lng": 10.007812058390755,
    "lat": 53.62953213568292
  },
  "address": "Lilienthalstr., 22335 Hamburg, DE",
  "uri": "/fs/5111fd35e4b0752e2e6d7219/details"
}
...
]
```
 
## Application Internals

### Dependencies & Libraries
This application relies on minimal Golang libraries in order to build the needed functionality. 

* This application uses [go dep](https://golang.github.io/dep/docs/introduction.html) to manage dependencies. 
* Gorilla [MUX](https://github.com/gorilla/mux) for routing and Gorilla [Handlers](https://github.com/gorilla/handlers) for specific http requests handlers (e.g. panic recovery)
* Google [go client](https://github.com/googlemaps/google-maps-services-go) for Maps Services
* A [go client](https://github.com/peppage/foursquarego) for Foursquare's API
  
### Application components

1) **places.go** : the main bootstrapping file of the application and where Dependencies gets Injected.

2) package **providers** : hosts different places providers. Currently *Google Places* and *Foursquare* are implemented. All providers (have to) implement the **Provider** interface. A strategy to allow adding new Providers and consuming them generically. 

3) package **handlers** : hosts different http handlers. First entry point for a user requests. The package host also different middlewares
    * A **loggingMiddleware** to log all requests
    * A **requestIdMiddleware** responsible to assigning a unique request id (correlation id) for every http request. For a better traceability
    * An **errorHandler** used as a place for centralized error handling. (The Go way is to return back all downstream errors to the caller. In this centralized handler we can judge how to handle different error types and thus, have default fallback scenarios.  

    The places.go (main) file, uses other handlers like a RecoveryHandler in order to recover the application from any Go "panic"(s).  

4) package **api** : hosts the webservice API resources definitions/models. 

5) package **config** : a basic package to load application configuration. Usually (especially in a microservice architecture) your service can be connected to a configuration service. In other setup(s)) config-maps/files can be mounted to your container and can be used for an application configuration (as an example, see Kubernetes'[configmaps](https://kubernetes.io/docs/tasks/configure-pod-container/configure-pod-configmap/)).

6) There is a simple Makefile in this repository to automate casual tasks (test, build..). However, it is better to hook a CI tool with this repo (out of this demo'scope)

Note on the implementation:

* Firing calls and getting places from different providers is done in parallel using Go routines. Please refer to the `handlers > handlers.go` 

## Local development

You can extend this code locally by either:

1) Getting the code and developing locally. Use `go get -u <repo>` 

You need a local [Golang setup](https://golang.org/doc/install). 
    
2) or you can simply mount the app as a volume in a golang container 

For example:

```sh
docker run -it --rm -p 8081:8081 --volume "$PWD":/go/src/app --workdir /go/src/app golang:1.11-alpine go run places.go [...other args]
```

## notes & future improvements

* In case you are using local Golang installation with locale paths for project source code: (must be in a src folder) https://github.com/golang/dep/issues/911
* [Improve]: add additional general handlers like Compress and CORS : straight forward implementation. Look in Gorilla read-made [handlers](https://github.com/gorilla/handlers)
* [Improve]: Extend tests for a higher code coverage (due to time constraints). The tests can already demonstrate how to test handlers, how to mock dependencies and how to imitate http calls. 

