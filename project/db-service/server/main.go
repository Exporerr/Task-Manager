package main

import (
	"context"

	logger "myproject/project/Logger"
	handlers "myproject/project/db-service/Handlers"
	databaseconnect "myproject/project/db-service/database_connect"
	"myproject/project/db-service/database_connect/service"

	"github.com/gorilla/mux"

	"net/http"
)

func main() {
	logger := logger.NewLogger()
	ctx := context.Background()
	url, erro := databaseconnect.LoadConfig("database.yml")
	if erro != nil {
		logger.Error.Fatalf("error to lodconfig: %v", erro)
	}

	pool, err := databaseconnect.NewPool(ctx, url.DatabaseURL)
	if err != nil {
		logger.Error.Fatalf("failed to connect to db: %v", err)
	}
	defer pool.Close()

	repo := databaseconnect.NewUserPool(pool, logger)
	logger.Info.Println("Repository Created")
	s := service.NewService(repo, logger)
	logger.Info.Println("Service Created")
	h := handlers.NewHandler(*s)
	logger.Info.Println("Handler Created")

	r := mux.NewRouter()
	r.HandleFunc("/tasks", h.Post).Methods("POST")
	r.HandleFunc("/tasks/{id}", h.GetTask).Methods("GET")
	r.HandleFunc("/tasks", h.AllTasks).Methods("GET")
	r.HandleFunc("/tasks/{id}", h.Patch).Methods("PATCH")
	r.HandleFunc("/tasks/{id}", h.Delete).Methods("DELETE")

	logger.Info.Println("Server started at :8081")
	if err := http.ListenAndServe(":8081", r); err != nil {
		logger.Error.Fatalf("server failed: %v", err)
	}
}
