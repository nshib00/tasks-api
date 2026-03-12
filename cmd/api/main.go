package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"tasks-api/internal/database"
	"tasks-api/internal/handlers"
)

func main() {
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	if dbUser == "" || dbPassword == "" || dbHost == "" || dbPort == "" || dbName == "" {
		log.Fatal("error: some of env variables are not set\n")
		os.Exit(1)
	}

	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", dbUser, dbPassword, dbHost, dbPort, dbName)
	db, err := database.Connect(dbURL)
	if err != nil {
		log.Fatal("database connection error: %v", err)
		return
	}
	defer db.Close()
	log.Println("connected to database successfully")

	repo := database.NewTaskRepository(db)
	handler := handlers.NewHandler(repo)

	mux := http.NewServeMux()
	mux.HandleFunc("/api/tasks", taskMethodHandler(handler.GetAllTasks, http.MethodGet))
	mux.HandleFunc("/api/tasks/create", taskMethodHandler(handler.CreateTask, http.MethodPost))
	mux.HandleFunc("/api/tasks/", taskMethodWithIDHandler(handler))
}

func taskMethodHandler(handlerFunc http.HandlerFunc, allowedMethod string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != allowedMethod {
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		}
		handlerFunc(w, r)
	}
}

func taskMethodWithIDHandler(handler *handlers.TasksHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handler.GetTaskByID(w, r)
		case http.MethodPut:
			handler.UpdateTask(w, r)
		case http.MethodDelete:
			handler.DeleteTask(w, r)
		default:
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		}
	}
}
