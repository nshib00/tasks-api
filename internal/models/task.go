package models

import (
	"time"
)


type Task struct {
	ID int `json:"id" db:"id"`
	Title string `json:"title" db:"title"`
	Description string `json:"description" db:"description"`
	Completed bool `json:"completed" db:"completed"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type CreateTaskInput string {
	Title string `json:"title"`
	Description string `json:"description"`
	Completed bool `json:"completed"`
}

type UpdateTaskInput string {
	Title *string `json:"title"`
	Description *string `json:"description"`
	Completed *bool `json:"completed"`
}