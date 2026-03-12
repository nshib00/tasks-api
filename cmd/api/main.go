package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"tasks-api/internal/database"
	"tasks-api/internal/handlers"
	"tasks-api/internal/middleware"
)

func getFromEnv(varName string) string {
	envVar := os.Getenv(varName)
	if envVar == "" {
		log.Fatalf("error: env variable does not exist or not set: %s", varName)
	}
	return envVar
}

func main() {
	dbUser := getFromEnv("POSTGRES_USER")
	dbPassword := getFromEnv("POSTGRES_PASSWORD")
	dbHost := getFromEnv("POSTGRES_HOST")
	dbPort := getFromEnv("POSTGRES_PORT")
	dbName := getFromEnv("POSTGRES_NAME")
	serverPort := getFromEnv("SERVER_PORT")

	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", dbUser, dbPassword, dbHost, dbPort, dbName)
	db, err := database.Connect(dbURL)
	if err != nil {
		log.Fatal("database connection error:", err)
	}
	defer db.Close()
	log.Println("connected to database successfully")

	repo := database.NewTaskRepository(db)
	handler := handlers.NewHandler(repo)

	mux := http.NewServeMux()
	mux.HandleFunc("/api/tasks", taskMethodHandler(handler.GetAllTasks, http.MethodGet))
	mux.HandleFunc("/api/tasks/create", taskMethodHandler(handler.CreateTask, http.MethodPost))
	mux.HandleFunc("/api/tasks/", taskMethodWithIDHandler(handler))

	loggedMux := middleware.LoggingMiddleware(mux)
	serverAddr := ":" + serverPort
	if err := http.ListenAndServe(serverAddr, loggedMux); err != nil {
		log.Fatal("server error: ", err)
	}
	log.Printf("server started on %s", serverAddr)
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
