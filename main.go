package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

func GetHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "text/plain")

	fmt.Fprintf(w, "Hello Echo")
}

func PostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)

	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain")

	fmt.Fprintf(w, "Received Post req with body: %s", string(body))

}
func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/get", GetHandler)
	mux.HandleFunc("/post", PostHandler)

	fmt.Println("Server Listening on http://localhost:8080")

	log.Fatal(http.ListenAndServe(":8080", mux))
}
