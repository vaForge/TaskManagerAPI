package store

import (
	"errors"
	"strings"

	"github.com/vaForge/TaskManagerAPI/models"
)

//Store holds task data in memory now
// This is a learning version
// Later replace with a database

type Store struct {
	tasks  []models.Task
	nextID int
}

// Initialise New Store
func NewStore() *Store {
	return &Store{
		tasks:  []models.Task{},
		nextID: 1,
	}
}

// GetAll returns all tasks
func (s *Store) GetAll() []models.Task {
	return s.tasks
}

// GetByID finds and returns one task
func (s *Store) GetByID(id int) (models.Task, bool) {
	for _, task := range s.tasks {
		if task.ID == id {
			return task, true
		}
	}
	return models.Task{}, false
}

// Create validates and adds a new Task
func (s *Store) Create(task models.Task) (models.Task, error) {

	task.Description = strings.TrimSpace(task.Description)

	if task.Description == "" {
		return models.Task{}, errors.New("description is required")
	}

	task.ID = s.nextID
	s.nextID++

	s.tasks = append(s.tasks, task)
	return task, nil
}

// Update replaces an existing task by ID
func (s *Store) Update(id int, updated models.Task) (models.Task, error) {

	index := -1

	for i, task := range s.tasks {
		if task.ID == id {
			index = i
			break
		}
	}

	if index == -1 {
		return models.Task{}, errors.New("task not found")

	}

	updated.Description = strings.TrimSpace(updated.Description)
	if updated.Description == "" {
		return models.Task{}, errors.New("description is required")
	}

	updated.ID = id
	s.tasks[index] = updated

	return updated, nil
}

// Delete remove a task by ID
func (s *Store) Delete(id int) error {

	index := -1

	for i, task := range s.tasks {
		if task.ID == id {
			index = i
			break
		}
	}

	if index == -1 {
		return errors.New("task not found")

	}

	s.tasks = append(s.tasks[:index], s.tasks[index+1:]...)
	return nil
}
