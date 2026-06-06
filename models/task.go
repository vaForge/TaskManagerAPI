package models

type Task struct {
	ID          int    `json:"id"`
	Description string `json:"description"`
	Status      bool   `json:"status"`
	Priority    int    `json:"priority"`
}
