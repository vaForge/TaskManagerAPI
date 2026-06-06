package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/vaForge/TaskManagerAPI/handlers"
	"github.com/vaForge/TaskManagerAPI/store"
)

func main() {
	// Create taskStore
	taskStore := store.NewStore()

	//Create taskhandlers
	taskhandler := handlers.NewTaskHandler(taskStore)

	mux := http.NewServeMux()

	//collection routes
	mux.HandleFunc("/tasks", taskhandler.TaskHandler)
	// Item routes
	mux.HandleFunc("/tasks/", taskhandler.TaskByIDHandler)
	// home route
	mux.HandleFunc("/", homeHandler)

	fmt.Println("Server Listening on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", mux))

}

func homeHandler(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	fmt.Fprintln(w, "Welcome to Task Manager API bros")
}
