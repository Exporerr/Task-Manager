package main

import (
	logger "myproject/project/Logger"
	"myproject/project/api-service/client"
	"myproject/project/api-service/handlers"
	"myproject/project/api-service/service"
	"myproject/project/middleware"
	"myproject/project/shared"

	"log"
	"net/http"

	"os"

	"github.com/gorilla/mux"
	"gopkg.in/yaml.v3"
)

func main() {
	logger := logger.NewLogger()
	cfg := shared.Config{}
	data, _ := os.ReadFile("config.yaml")
	yaml.Unmarshal(data, &cfg)
	client := client.NewClient(cfg.DBService.URL, *logger)
	service := service.NewService(client, logger)
	handler := handlers.NewHandler(*service, logger)

	r := mux.NewRouter()
	r.Use(middleware.LoggingMiddlware)

	r.HandleFunc("/tasks", handler.Post).Methods("POST")
	r.HandleFunc("/tasks/{id}", handler.Get).Methods("GET")
	r.HandleFunc("/tasks", handler.GetAll).Methods("GET")
	r.HandleFunc("/tasks/{id}", handler.Update).Methods("PATCH")
	r.HandleFunc("/tasks/{id}", handler.Delete).Methods("DELETE")

	log.Println("Server started at :8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
