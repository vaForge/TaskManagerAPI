package models

type Task struct {
	ID          int    `json:"id"`
	Description string `json:"description"`
	Status      bool   `json:"status"`
	Priority    int    `json:"priority"`
}

// Task Patch is used for PATCH requests

// Pointer fields let us tell the difference between
// - field not provided at all ( nil )
// - field provided with zero values (eg : false, 0, "")

// That matters more parital updates
type TaskPatch struct {
	Description *string `json:"description,omitempty"`
	Status      *bool   `json:"status,omitempty"`
	Priority    *int    `json:"priority,omitempty"`
}
