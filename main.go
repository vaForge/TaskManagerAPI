package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/vaForge/TaskManagerAPI/models"
)

// tasks stores the in-memory task list.
// This is temporary storage only.
// When the server restarts, the data is lost.
var tasks []models.Task

// idCounter generates unique task IDs.
var idCounter = 1

func main() {
	mux := http.NewServeMux()

	// Collection route:
	// GET  /tasks      -> list all tasks
	// POST /tasks      -> create a new task
	mux.HandleFunc("/tasks", tasksHandler)

	// Item route:
	// GET    /tasks/{id}    -> get one task
	// PUT    /tasks/{id}    -> update one task
	// DELETE /tasks/{id}    -> delete one task
	mux.HandleFunc("/tasks/", taskByIDHandler)

	// Simple home page to test the server quickly.
	mux.HandleFunc("/", homeHandler)

	fmt.Println("Server listening on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}

// homeHandler is just a small welcome endpoint.
func homeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	fmt.Fprintln(w, "Welcome to the Task Manager API!")
}

// tasksHandler handles the collection endpoint: /tasks
func tasksHandler(w http.ResponseWriter, r *http.Request) {
	// Make sure this handler is only used for the exact /tasks path.
	if r.URL.Path != "/tasks" {
		http.NotFound(w, r)
		return
	}

	switch r.Method {
	case http.MethodGet:
		getTasks(w, r)
	case http.MethodPost:
		createTask(w, r)
	default:
		w.Header().Set("Allow", "GET, POST")
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

// taskByIDHandler handles /tasks/{id}
func taskByIDHandler(w http.ResponseWriter, r *http.Request) {
	// Expected path format:
	// /tasks/1
	// /tasks/2
	//
	// We trim the "/tasks/" prefix and parse the remaining part as an integer.
	idStr := strings.TrimPrefix(r.URL.Path, "/tasks/")

	// If the path is just "/tasks/" or malformed, reject it.
	if idStr == "" || strings.Contains(idStr, "/") {
		http.Error(w, "invalid task id", http.StatusBadRequest)
		return
	}

	taskID, err := strconv.Atoi(idStr)
	if err != nil || taskID <= 0 {
		http.Error(w, "invalid task id", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		getTaskByID(w, r, taskID)
	case http.MethodPut:
		updateTaskByID(w, r, taskID)
	case http.MethodDelete:
		deleteTaskByID(w, r, taskID)
	default:
		w.Header().Set("Allow", "GET, PUT, DELETE")
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

// getTasks returns all tasks as a JSON array.
func getTasks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Encode the full slice into JSON.
	// Example output:
	// [
	//   {"id":1,"description":"Learn Go","status":false,"priority":1}
	// ]
	if err := json.NewEncoder(w).Encode(tasks); err != nil {
		http.Error(w, "failed to encode tasks", http.StatusInternalServerError)
		return
	}
}

// createTask reads a JSON body, assigns an ID, and stores the task.
func createTask(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var newTask models.Task

	// Decode JSON directly from request body.
	if err := json.NewDecoder(r.Body).Decode(&newTask); err != nil {
		http.Error(w, "invalid JSON body", http.StatusBadRequest)
		return
	}

	// Basic validation.
	if strings.TrimSpace(newTask.Description) == "" {
		http.Error(w, "description is required", http.StatusBadRequest)
		return
	}

	// Server controls the ID.
	newTask.ID = idCounter
	idCounter++

	// Save in memory.
	tasks = append(tasks, newTask)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	// Return the created task.
	if err := json.NewEncoder(w).Encode(newTask); err != nil {
		http.Error(w, "failed to encode created task", http.StatusInternalServerError)
		return
	}
}

// getTaskByID finds one task and returns it as JSON.
func getTaskByID(w http.ResponseWriter, r *http.Request, taskID int) {
	task, found := findTaskByID(taskID)
	if !found {
		http.Error(w, "task not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(task); err != nil {
		http.Error(w, "failed to encode task", http.StatusInternalServerError)
		return
	}
}

// updateTaskByID updates one task by ID.
func updateTaskByID(w http.ResponseWriter, r *http.Request, taskID int) {
	defer r.Body.Close()

	index, found := findTaskIndexByID(taskID)
	if !found {
		http.Error(w, "task not found", http.StatusNotFound)
		return
	}

	// Read the updated data from the request body.
	var updatedTask models.Task
	if err := json.NewDecoder(r.Body).Decode(&updatedTask); err != nil {
		http.Error(w, "invalid JSON body", http.StatusBadRequest)
		return
	}

	// Keep the original ID.
	updatedTask.ID = taskID

	// Replace the old task with the new one.
	tasks[index] = updatedTask

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(updatedTask); err != nil {
		http.Error(w, "failed to encode updated task", http.StatusInternalServerError)
		return
	}
}

// deleteTaskByID removes a task from the slice.
func deleteTaskByID(w http.ResponseWriter, r *http.Request, taskID int) {
	index, found := findTaskIndexByID(taskID)
	if !found {
		http.Error(w, "task not found", http.StatusNotFound)
		return
	}

	// Remove the item at index from the slice.
	tasks = append(tasks[:index], tasks[index+1:]...)

	// 204 No Content means deletion succeeded and no response body is needed.
	w.WriteHeader(http.StatusNoContent)
}

// findTaskByID returns the task and whether it exists.
func findTaskByID(taskID int) (models.Task, bool) {
	for _, task := range tasks {
		if task.ID == taskID {
			return task, true
		}
	}
	return models.Task{}, false
}

// findTaskIndexByID returns the index of the task and whether it exists.
func findTaskIndexByID(taskID int) (int, bool) {
	for i, task := range tasks {
		if task.ID == taskID {
			return i, true
		}
	}
	return -1, false
}
