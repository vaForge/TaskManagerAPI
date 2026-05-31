package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

type Task struct {
	ID          int    `json:"id"`
	Description string `json:"description"`
	Status      bool   `json:"status"`
	Priority    int    `json:"priority"`
}

var tasks []Task

func GetHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	for _, task := range tasks {
		taskJson, err := json.Marshal(task)
		if err != nil {
			http.Error(w, "Failed to marshal task", http.StatusInternalServerError)
			return
		}
		// w.Write(taskJson)
		// w.Write([]byte("\n"))
		fmt.Fprintf(w, "%s\n", string(taskJson))
	}
}

func PostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed) //405 Method Not Allowed
		return
	}

	body, err := io.ReadAll(r.Body)
	defer r.Body.Close()

	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusInternalServerError) //500 Internal Server Error
		return
	}

	var newTask Task
	err = json.Unmarshal(body, &newTask)
	if err != nil {
		http.Error(w, "Failed to parse JSON", http.StatusBadRequest) //400 Bad Request
		return
	}
	tasks = append(tasks, newTask)
	w.WriteHeader(http.StatusCreated) //201 Created
	w.Write([]byte("Task Created Successfully\n"))

}
func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/get", GetHandler)
	mux.HandleFunc("/post", PostHandler)

	fmt.Println("Server Listening on http://localhost:8080")

	log.Fatal(http.ListenAndServe(":8080", mux))
}
