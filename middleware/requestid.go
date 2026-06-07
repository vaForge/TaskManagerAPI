package middleware

// Each request gets a unique ID so you can trace logs later.

// That becomes very useful when debugging a failing request.

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"
)

const requestIDHeader = "X-Request-Id"

// RequestID ensures every request has a request ID header.

func RequestID(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		id := r.Header.Get(requestIDHeader)

		if id == "" {
			id = newRequestID()
		}

		w.Header().Set(requestIDHeader, id)
		next.ServeHTTP(w, r)
	})

}

func newRequestID() string {
	var b [16]byte
	_, _ = rand.Read(b[:])

	dst := make([]byte, 32)
	hex.Encode(dst, b[:])
	return string(dst)
}
