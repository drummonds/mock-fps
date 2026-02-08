package jsonapi

import (
	"net/http"
	"strings"
)

const ContentType = "application/vnd.api+json"

// EnforceContentType is middleware that validates Content-Type on requests with bodies.
func EnforceContentType(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost || r.Method == http.MethodPatch || r.Method == http.MethodPut {
			ct := r.Header.Get("Content-Type")
			if ct != "" && !strings.HasPrefix(ct, ContentType) && !strings.HasPrefix(ct, "application/json") {
				WriteError(w, http.StatusUnsupportedMediaType, "Unsupported Media Type",
					"Content-Type must be application/vnd.api+json")
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}
