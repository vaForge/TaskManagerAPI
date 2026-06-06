package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/vaForge/TaskManagerAPI/models"
	"github.com/vaForge/TaskManagerAPI/store"
)

// TaskHandler owns the HTTP layer
// it talks to the store but does not store the data itself
// This separation is useful as project grows

type TaskHandler struct {
	Store *store.Store
}

func NewTaskHandler(s *store.Store) *TaskHandler {
	return &TaskHandler{Store: s}
}

// TaskHandler handles /tasks
// Get /tasks -> all tasks
// Post /tasks -> create task

func (h *TaskHandler) TaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/tasks" {
		http.NotFound(w, r)
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.getTasks(w, r)
	case http.MethodPost:
		h.createTask(w, r)
	default:
		w.Header().Set("Allow", "GET, POST")
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

// TaskByIDHandler handles /tasks/{id}
// Get     /tasks/{id}
// PUT 		/tasks/{id}
// Delete 	/tasks/{id}

func (h *TaskHandler) TaskByIDHandler(w http.ResponseWriter, r *http.Request) {
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
		h.getTaskByID(w, r, taskID)
	case http.MethodPut:
		h.updateTaskByID(w, r, taskID)
	case http.MethodDelete:
		h.deleteTaskByID(w, r, taskID)
	default:
		w.Header().Set("Allow", "GET, PUT, DELETE")
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *TaskHandler) getTasks(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(h.Store.GetAll()); err != nil {
		http.Error(w, "failed to encode tasks", http.StatusInternalServerError)
		return
	}
}

func (h *TaskHandler) createTask(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var newTask models.Task

	// Decode JSON directly from request body.
	if err := json.NewDecoder(r.Body).Decode(&newTask); err != nil {
		http.Error(w, "invalid JSON body", http.StatusBadRequest)
		return
	}

	createdTask, err := h.Store.Create(newTask)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(createdTask); err != nil {
		http.Error(w, "failed to encode task", http.StatusInternalServerError)
		return
	}
}

func (h *TaskHandler) getTaskByID(w http.ResponseWriter, r *http.Request, taskID int) {

	task, found := h.Store.GetByID(taskID)

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

func (h *TaskHandler) updateTaskByID(w http.ResponseWriter, r *http.Request, taskID int) {

	defer r.Body.Close()

	var updatedTask models.Task

	if err := json.NewDecoder(r.Body).Decode(&updatedTask); err != nil {
		http.Error(w, "invalid JSON body", http.StatusBadRequest)
		return
	}

	result, err := h.Store.Update(taskID, updatedTask)

	if err != nil {
		if err.Error() == "task not found" {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(result); err != nil {
		http.Error(w, "failed to encode updated task", http.StatusInternalServerError)
		return
	}
}

func (h *TaskHandler) deleteTaskByID(w http.ResponseWriter, r *http.Request, taskID int) {

	if err := h.Store.Delete(taskID); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
