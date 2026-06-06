package middleware

import (
	"log"
	"net/http"
	"time"
)

type statusRecorder struct {
	http.ResponseWriter
	status int
	size   int
}

// WriteHeader stoers the status code before sending it to the client

func (sr *statusRecorder) WriteHeader(code int) {
	sr.status = code
	sr.ResponseWriter.WriteHeader(code)
}

// Write stores the body size and sets default stauts to 200 if needed.
func (sr *statusRecorder) Write(b []byte) (int, error) {
	if sr.status == 0 {
		sr.status = http.StatusOK

	}

	n, err := sr.ResponseWriter.Write(b)
	sr.size = n

	return n, err
}

// Logging prints one line for every request
func Logging(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		start := time.Now()

		rec := &statusRecorder{
			ResponseWriter: w,
		}

		next.ServeHTTP(rec, r)

		log.Printf("%s %s %d %s",
			r.Method,
			r.URL.Path,
			rec.status,
			time.Since(start),
		)
	})
}
