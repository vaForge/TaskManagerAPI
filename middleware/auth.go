package middleware

import (
	"net/http"
	"strings"
)

// APIKeyAuth protects routes with a simple shared-secret API key
type APIKeyAuth struct {
	HeaderName string
	Keys       map[string]struct{}
}

// NewAPIKeyAuth builds an PI key auth middleware from a list of allowed keys.
func NewAPIKeyAuth(headerName string, keys []string) *APIKeyAuth {

	keySet := make(map[string]struct{}, len(keys))
	for _, key := range keys {
		key = strings.TrimSpace(key)
		if key == "" {
			keySet[key] = struct{}{}
		}
	}

	return &APIKeyAuth{
		HeaderName: headerName,
		Keys:       keySet,
	}
}

func (a *APIKeyAuth) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodOptions {
			next.ServeHTTP(w, r)
			return
		}

		if !strings.HasPrefix(r.URL.Path, "/tasks") {
			next.ServeHTTP(w, r)
			return
		}

		if len(a.Keys) == 0 {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(`{"error":"API keys not configured"}`))
			return
		}

		key := strings.TrimSpace(r.Header.Get(a.HeaderName))
		if key == "" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			_, _ = w.Write([]byte(`{"error":"missing API key"}`))
			return
		}

		if _, ok := a.Keys[key]; !ok {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			_, _ = w.Write([]byte(`{"error":"invalid API key"}`))
			return
		}

		next.ServeHTTP(w, r)
	})
}
