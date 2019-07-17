package handlers

import (
	"github.com/codeselim/go-webservice-places-provider/log"
	"net/http"
)

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// log and continue
		logger := log.GetLoggerWithContext(r.Context())
		logger.Info(r.Method, r.RequestURI, r.RemoteAddr, r.Host)
		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)
	})
}
