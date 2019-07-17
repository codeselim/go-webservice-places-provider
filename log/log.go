package log

import (
	"context"
	"github.com/codeselim/go-webservice-places-provider/config"
	"github.com/sirupsen/logrus"
	"os"
)

//// ContextKey is used for context.Context value. The value requires a key that is not primitive type.
//type ContextKey string // can be unexported
//
//// ContextKeyRequestID is the ContextKey for RequestID
//const ContextKeyRequestID ContextKey = "requestID" // can be unexported

var log = logrus.New()

func init() {
	// Log as JSON instead of the default ASCII formatter.
	log.SetFormatter(&logrus.JSONFormatter{})
	// Output to stdout instead of the default stderr
	log.SetOutput(os.Stdout)
	// Only log the warning severity or above.
	log.SetLevel(logrus.InfoLevel) //this can be extended in the future to be supplied via application flag
}

func GetLoggerWithContext(ctx context.Context) *logrus.Entry {
	return log.WithFields(logrus.Fields{
		"requestID": getRequestID(ctx),
	})

}

func GetLogger() *logrus.Entry {
	return log.WithFields(logrus.Fields{
		"requestID": "",
	})
}

func getRequestID(ctx context.Context) string {
	reqID := ctx.Value(config.ContextKeyRequestID)

	if ret, ok := reqID.(string); ok {
		return ret
	}

	return ""
}
