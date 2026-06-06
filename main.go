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

var tasks []models.Task

var idCounter int = 1

func main() {
	mux := http.NewServeMux()

	// collection route :
	// Get /tasks
	// POST /tasks
	mux.HandleFunc("/tasks", tasksHandler)

	//Item route:
	//Get /tasks/{id}
	//PUT /tasks/{id}
	//DELETE /tasks/{id}
	mux.HandleFunc("/task/", taskByIDHandler)

	//simple home page tester
	mux.HandleFunc("/", homeHandler)

	fmt.Println("Server Listening on http://localhost:8080")

	log.Fatal(http.ListenAndServe(":8080", mux))
}

// A simple welcome endpoint to test server
func homeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charsert=utf-8")
	fmt.Fprintln(w, "Welcome to the Task Manager API!")
}

func tasksHandler(w http.ResponseWriter, r *http.Request) {
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
		//Any method other than GET or POST is not allowed on this end point
		w.Header().Set("Allow", "GET, POST")
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func taskByIDHandler(w http.ResponseWriter, r *http.Request) {
	//Expected : /task/1 or /tasks/2
	//We trim "/tasks/" prefix and parse the remaining part as int

	idStr := strings.TrimPrefix(r.URL.Path, "/tasks/")

	if idStr == "" || strings.Contains(idStr, "/") {
		http.Error(w, "invalid task id", http.StatusBadRequest)
		return
	}

	taskID, err := strconv.Atoi(idStr)
	if err != nil {
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
		w.Header().Set("Allow", "GET,PUT, DELETE")
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}

}

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

func updateTaskByID(w http.ResponseWriter, r *http.Request, taskID int) {
	defer r.Body.Close()

	index, found := findTaskIndexByID(taskID)

	if !found {
		http.Error(w, "task not found", http.StatusNotFound)
		return
	}

	var updatedTask models.Task
	if err := json.NewDecoder(r.Body).Decode(&updatedTask); err != nil {
		http.Error(w, "invalid JSON format", http.StatusBadRequest)
		return
	}

	updatedTask.ID = taskID
	tasks[index] = updatedTask

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(updatedTask); err != nil {
		http.Error(w, "failed to encode updated task", http.StatusInternalServerError)
		return
	}

}
func deleteTaskByID(w http.ResponseWriter, r *http.Request, taskID int) {

	index, found := findTaskIndexByID(taskID)

	if !found {
		http.Error(w, "task not found", http.StatusNotFound)
		return
	}

	tasks = append(tasks[:index], tasks[index+1:]...)

	//204 Content means deletion succeeded and no response body is needed

	w.WriteHeader(http.StatusNoContent)
}

func getTasks(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	//Encode the whole slice at once
	//This returns a valid JSON like
	// [
	// 	{..},
	// 	{...}
	// ]

	if err := json.NewEncoder(w).Encode(tasks); err != nil {
		http.Error(w, "failed to encode tasks", http.StatusInternalServerError)
		return
	}
}

func createTask(w http.ResponseWriter, r *http.Request) {

	defer r.Body.Close()

	//Decode JSON directly from request body.
	//Cleaner than io.ReadAll + json.Unmarshal for APIs
	var newTask models.Task

	if err := json.NewDecoder(r.Body).Decode(&newTask); err != nil {

		http.Error(w, "invalid JSON body", http.StatusBadRequest) //400 Bad Request
		return
	}

	newTask.ID = idCounter
	idCounter++
	//save tasks
	tasks = append(tasks, newTask)

	//Return a proper JSON with 201 Created
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated) //201 Created

	//Send the created task back to client
	//This helps client see the generated ID immediately
	if err := json.NewEncoder(w).Encode(newTask); err != nil {
		http.Error(w, "failed to encode created task", http.StatusInternalServerError) //500
		return
	}

}

// FindtaskbyId returns the task and whether it exists
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
