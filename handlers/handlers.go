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
		writeError(w, http.StatusBadRequest, "invalid task id")
		return
	}

	taskID, err := strconv.Atoi(idStr)
	if err != nil || taskID <= 0 {
		writeError(w, http.StatusBadRequest, "invalid task id")
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.getTaskByID(w, r, taskID)
	case http.MethodPut:
		h.updateTaskByID(w, r, taskID)
	case http.MethodPatch:
		h.patchTaskByID(w, r, taskID)
	case http.MethodDelete:
		h.deleteTaskByID(w, r, taskID)
	default:
		w.Header().Set("Allow", "GET, PUT, DELETE, PATCH")
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *TaskHandler) getTasks(w http.ResponseWriter, r *http.Request) {

	writeJSON(w, http.StatusOK, h.Store.GetAll())
}

func (h *TaskHandler) createTask(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	if !requireJSONContentType(w, r) {
		return
	}

	var newTask models.Task

	// Decode JSON directly from request body.
	if err := json.NewDecoder(r.Body).Decode(&newTask); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}

	createdTask, err := h.Store.Create(newTask)

	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, createdTask)
}

func (h *TaskHandler) getTaskByID(w http.ResponseWriter, r *http.Request, taskID int) {

	task, found := h.Store.GetByID(taskID)

	if !found {
		writeError(w, http.StatusNotFound, "task not found")
		return
	}

	writeJSON(w, http.StatusOK, task)
}

func (h *TaskHandler) updateTaskByID(w http.ResponseWriter, r *http.Request, taskID int) {

	defer r.Body.Close()

	if !requireJSONContentType(w, r) {
		return
	}

	var updatedTask models.Task

	if err := json.NewDecoder(r.Body).Decode(&updatedTask); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}

	result, err := h.Store.Update(taskID, updatedTask)

	if err != nil {
		if err.Error() == "task not found" {
			writeError(w, http.StatusNotFound, err.Error())
			return
		}
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, result)
}

func (h *TaskHandler) patchTaskByID(w http.ResponseWriter, r *http.Request, taskID int) {

	defer r.Body.Close()

	if !requireJSONContentType(w, r) {
		return
	}

	var patch models.TaskPatch

	if err := json.NewDecoder(r.Body).Decode(&patch); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}

	result, err := h.Store.Patch(taskID, patch)

	if err != nil {
		if err.Error() == "task not found" {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, result)
}

func (h *TaskHandler) deleteTaskByID(w http.ResponseWriter, r *http.Request, taskID int) {

	if err := h.Store.Delete(taskID); err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
