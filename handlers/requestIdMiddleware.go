package handlers

import (
	"context"
	"github.com/codeselim/go-webservice-places-provider/config"
	"github.com/codeselim/go-webservice-places-provider/log"
	"github.com/google/uuid"
	"net/http"
)

// AttachRequestID will attach a brand new request ID to a http request
func AssignRequestID(ctx context.Context) context.Context {

	reqID := uuid.New()

	return context.WithValue(ctx, config.ContextKeyRequestID, reqID.String())
}

// GetRequestID will get reqID from a http request and return it as a string
func GetRequestID(ctx context.Context) string {

	reqID := ctx.Value(config.ContextKeyRequestID)

	if ret, ok := reqID.(string); ok {
		return ret
	}

	return ""
}

func RequestIdMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		r = r.WithContext(AssignRequestID(ctx))

		logger := log.GetLoggerWithContext(r.Context())

		logger.Info("Incomming http request")

		next.ServeHTTP(w, r)

		logger.Info("Finished handling http request")
	})
}
