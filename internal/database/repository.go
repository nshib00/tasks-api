package database

import (
	"database/sql"
	"errors"
	"tasks-api/internal/models"
	"time"
	appErrors "tasks-api/internal/errors"
	"github.com/jmoiron/sqlx"
)

type TasksRepository struct {
	db *sqlx.DB
}

func NewTaskRepository(db *sqlx.DB) *TasksRepository {
	return &TasksRepository{db: db}
}

func (repo *TasksRepository) Create(input models.CreateTaskInput) (*models.Task, error) {
	var task models.Task

	query := `
		INSERT INTO tasks (title, description, completed, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, title, description, completed, created_at;
	`
	now := time.Now()

	err := repo.db.QueryRowx(
		query, input.Title, input.Description, input.Completed, now, now,
	).StructScan(&task)
	if err != nil {
		return nil, err
	}
	return &task, nil
}

func (repo *TasksRepository) GetAll() ([]models.Task, error) {
	var tasks []models.Task

	query := `
		SELECT title, description, completed, created_at, updated_at
		FROM tasks
		ORDER BY updated_at DESC
	`
	err := repo.db.Select(&tasks, query)
	if err != nil {
		return nil, err
	}
	return tasks, nil
}

func (repo *TasksRepository) GetByID(id int) (*models.Task, error) {
	var task *models.Task

	query := `
		SELECT description, completed, created_at, updated_at
		FROM tasks
		WHERE id = $1
	`
	err := repo.db.Get(&task, query, id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, errors.New("task not found")
	}
	if err != nil {
		return nil, err
	}
	return &task, nil
}

func (repo *TasksRepository) Update(id int, input models.UpdateTaskInput) (*models.Task, error) {
	task, err := repo.GetByID(id)
	if err != nil {
		return nil, appErrors.TaskNotFound
	}

	if input.Title != nil {
		task.Title = *input.Title
	}
	if input.Description != nil {
		task.Description = *input.Description
	}
	if input.Completed != nil {
		task.Completed = *input.Completed
	}
	task.UpdatedAt := time.Now()

	var updatedTask *models.Task
	query := `
		UPDATE tasks
		SET title=$1, description=$2, completed=$3, updated_at=$4
		WHERE id=$5
		RETURNING id, title, description, completed, updated_at
	`
	err := repo.db.QueryRowx(
		query, task.Title, task.Description, task.Completed, task.UpdatedAt, id
	).StructScan(&updatedTask)
	if err != nil {
		return nil, err
	}
	return &task, nil
}

func (repo *TasksRepository) Delete(id int) error {
	query = `DELETE FROM tasks WHERE id=$1`
	result, err := repo.db.Exec(query, id)
	affectedRows := result.RowsAffected()

	if affectedRows == 0 {
		return appErrors.TaskNotFound
	} 
	if err != nil {
		return err
	}
	return nil
}
