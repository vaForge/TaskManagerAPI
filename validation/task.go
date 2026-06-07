package validation

import (
	"errors"
	"net/http"
	"strings"

	"github.com/vaForge/TaskManagerAPI/models"
)

// Validation checks the full task ovjects

// This func for POST and PUT reqs
func ValidateTask(task models.Task) error {

	if strings.TrimSpace(task.Description) == "" {
		return errors.New("description is required")
	}

	//Simple learning rule :
	//Priority must be between 1 and 6
	if task.Priority < 1 || task.Priority > 5 {
		return errors.New("priority must be between 1 and 6")
	}

	return nil
}

//ValidateTaskpatch checks only the fields the were actually sent
// Used for PATCH requests

func ValidateTaskPatch(patch models.TaskPatch) error {

	if patch.Description != nil && strings.TrimSpace(*patch.Description) == "" {
		return errors.New("description cannot be empty")
	}

	if patch.Priority != nil && (*patch.Priority < 1 || *patch.Priority > 5) {
		return errors.New("Priority must be between 1 and 6")
	}

	return nil
}

// IsJSONContentType checks whether the request body is JSON
// This is useful for POST PUT PATCH endpoints
func IsJSONContentType(r *http.Request) bool {
	contentType := r.Header.Get("Content-Type")
	return strings.HasPrefix(contentType, "application/json")
}
