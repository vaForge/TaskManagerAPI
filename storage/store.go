package main

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"
)

type Task struct{
	ID string `json:"id"`
	Title string `json:"title"`
	Done bool `json:"done"`
	CreatedAt time.Time `json:"created_at"`
}

type TaskStore struct {
	mu sync.RWMutex
	tasks map[string]Task
	filePath string
}

func NewTaskStore(filePath string)( *TaskStore ,error){
	s := &TaskStore{
		tasks : make(map[string]Task),
		filePath: filePath,
	}
	if err := s.load(); err != nil{
		return nil,fmt.Errorf("loading tasks: %w",err)
	}
	return s,nil
}	

func (s *TaskStore) load() error {
	data,err := os.ReadFile(s.filePath)
	if os.IsNotExist(err){
		return nil // fresh start
	}
	if err != nil {
		return err
	}

	var list []Task
	if err := json.Unmarshal(data,&list); err != nil {
		return err
	}
	for _,t := range list {
		s.tasks[t.ID] = t
	}
	return nil
}

func (s *TaskStore) save() error {
	list :=  make([]Task,0,len(s.tasks))
	for _,t := range s.tasks {
		list = append(list,t)
	}
	data,err := json.MarshalIndent(list," "," ")	
	if err != nil {
		return err
	}
	tmp := s.filePath + ".json"
	if err := os.WriteFile(tmp,data,0644); err != nil {
		return err
	}
	return os.Rename(tmp,s.filePath)
}

func newID() string {
	return fmt.Sprintf("%d",time.Now().UnixNano())
}
// All() - show all tasks , Add() - adds new tasks , Toggle() - mark as done 
// Delete() - delete task by ID