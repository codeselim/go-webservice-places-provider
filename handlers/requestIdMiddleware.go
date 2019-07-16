package handlers

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"log"
	"net/http"
)

// ContextKey is used for context.Context value. The value requires a key that is not primitive type.
type ContextKey string // can be unexported

// ContextKeyRequestID is the ContextKey for RequestID
const ContextKeyRequestID ContextKey = "requestID" // can be unexported

// AttachRequestID will attach a brand new request ID to a http request
func AssignRequestID(ctx context.Context) context.Context {

	reqID := uuid.New()

	return context.WithValue(ctx, ContextKeyRequestID, reqID.String())
}

// GetRequestID will get reqID from a http request and return it as a string
func GetRequestID(ctx context.Context) string {

	reqID := ctx.Value(ContextKeyRequestID)

	if ret, ok := reqID.(string); ok {
		return ret
	}

	return ""
}

func RequestIdMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		r = r.WithContext(AssignRequestID(ctx))

		log.Printf("Incomming request %s %s %s %s", r.Method, r.RequestURI, r.RemoteAddr, GetRequestID(r.Context()))

		next.ServeHTTP(w, r)
		fmt.Print(GetRequestID(ctx))
		log.Printf("Finished handling http req. %s", GetRequestID(r.Context()))
	})
}
