package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"
	"tasks-api/internal/database"
	appErrors "tasks-api/internal/errors"
	"tasks-api/internal/models"
)

type TasksHandler struct {
	repo *database.TasksRepository
}

func NewHandler(repo *database.TasksRepository) *TasksHandler {
	return &TasksHandler{
		repo: repo,
	}
}

func respondWithJSON(w http.ResponseWriter, statusCode int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		log.Printf("JSON encoding error: %v", err)
	}
}

func respondWithError(w http.ResponseWriter, statusCode int, message string) {
	respondWithJSON(w, statusCode, map[string]string{"error": message})
}

func (handler *TasksHandler) GetAllTasks(w http.ResponseWriter, _ *http.Request) {
	tasks, err := handler.repo.GetAll()
	if err != nil {
		log.Printf("error in TasksHandler.GetAllTasks: %v", err)
		respondWithError(w, http.StatusInternalServerError, "task receiving error")
		return
	}
	respondWithJSON(w, http.StatusOK, tasks)
}

func (handler *TasksHandler) GetTaskByID(w http.ResponseWriter, r *http.Request) {
	pathWithoutPrefix := strings.TrimPrefix(r.URL.Path, "/api/tasks/")
	pathParts := strings.Split(pathWithoutPrefix, "/")
	taskIDStr := pathParts[0]

	taskID, err := strconv.Atoi(taskIDStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "incorrect task ID")
		return
	}

	task, err := handler.repo.GetByID(taskID)
	if err != nil {
		switch {
		case errors.Is(err, appErrors.TaskNotFound):
			respondWithError(w, http.StatusNotFound, err.Error())
		default:
			log.Printf("Error in TaskHandler.GetByID: %v", err)
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
	}
	respondWithJSON(w, http.StatusOK, task)
}

func (handler *TasksHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
	var taskData models.CreateTaskInput

	if err := json.NewDecoder(r.Body).Decode(&taskData); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid task data")
		return
	}
	if strings.TrimSpace(taskData.Title) == "" {
		respondWithError(w, http.StatusBadRequest, "task title cannot be empty")
		return
	}

	task, err := handler.repo.Create(taskData)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error on task creating")
		log.Printf("Error in in TaskHandler.GetByID: %v", err)
		return
	}

	respondWithJSON(w, http.StatusOK, task)
}

func (handler *TasksHandler) UpdateTask(w http.ResponseWriter, r *http.Request) {
	pathParts := strings.Split(strings.TrimPrefix(r.URL.Path, "/api/tasks/"), "/")
	taskIDStr := pathParts[0]

	var taskData models.UpdateTaskInput
	if err := json.NewDecoder(r.Body).Decode(&taskData); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid task data")
		return
	}
	if taskData.Title == nil || strings.TrimSpace(*taskData.Title) == "" {
		respondWithError(w, http.StatusBadRequest, "task title cannot be empty")
		return
	}

	taskID, err := strconv.Atoi(taskIDStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "incorrect task ID")
		return
	}

	task, err := handler.repo.Update(taskID, taskData)
	if err != nil {
		if errors.Is(err, appErrors.TaskNotFound) {
			respondWithError(w, http.StatusNotFound, "task not found")
		} else {
			log.Printf("error on task updating: %v", err)
			respondWithError(w, http.StatusInternalServerError, "error on task updating")
		}
		return
	}

	respondWithJSON(w, http.StatusCreated, task)
}

func (handler *TasksHandler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	pathParts := strings.Split(strings.TrimPrefix(r.URL.Path, "/api/tasks/"), "/")
	taskIDStr := pathParts[0]

	taskID, err := strconv.Atoi(taskIDStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid task ID")
		return
	}

	err = handler.repo.Delete(taskID)
	if err != nil {
		if errors.Is(err, appErrors.TaskNotFound) {
			respondWithError(w, http.StatusNotFound, "task not found")
		} else {
			respondWithError(w, http.StatusInternalServerError, "error on task deleting")
		}
		return
	}

	respondWithJSON(w, http.StatusNoContent, nil)
}
