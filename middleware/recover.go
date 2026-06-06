package middleware

import (
	"encoding/json"
	"log"
	"net/http"
)

// Recover catches panics and turns them into a safe 500 response

// Without this a panic inside a handler can crash the whole server

func Recover(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		defer func() {

			if rec := recover(); rec != nil {
				log.Printf("panic recovered %v", rec)

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)

				//keep the response simple and consistent

				_ = json.NewEncoder(w).Encode(map[string]string{
					"error": "internal server error",
				})
			}
		}()

		next.ServeHTTP(w, r)
	})

}
