package validator

import (
	"main/src/loggingMiddleware"
	"net/http"
)

func Main(lrw *loggingMiddleware.LoggingResponseWriter) bool {

	err := validateSchema(lrw)
	if err != nil {
		http.Error(lrw.ResponseW, err.Error(), http.StatusInternalServerError)
		return false
	}

	return true
}
