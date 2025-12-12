package main

import (
	"fmt"
	"log"
	"net/http"
	"os" 

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go-microservice/metrics"

	"go-microservice/handlers"
	"go-microservice/services"
	"go-microservice/utils"

	"github.com/gorilla/mux"
)

func main() {

	minioEndpoint := os.Getenv("MINIO_ENDPOINT")
	if minioEndpoint == "" {
		minioEndpoint = "localhost:9000"
	}

	minioAccessKey := os.Getenv("MINIO_ROOT_USER")
	if minioAccessKey == "" {
		minioAccessKey = "minioadmin"
	}

	minioSecretKey := os.Getenv("MINIO_ROOT_PASSWORD")
	if minioSecretKey == "" {
		minioSecretKey = "minioadmin"
	}

	minioBucket := "user-files"

	// Запускаем логгер
	go utils.StartLogger()

	// 1. Инициализация слоев (Dependency Injection)
	userService := services.NewUserService()
	// Передаем полученные настройки в сервис
	integrationService := services.NewIntegrationService(minioEndpoint, minioAccessKey, minioSecretKey, minioBucket)

	userHandler := handlers.NewUserHandler(userService)
	integrationHandler := handlers.NewIntegrationHandler(integrationService)

	// 2. Роутер
	r := mux.NewRouter()
	r.Handle("/metrics", promhttp.Handler())

	// 3. Регистрация маршрутов
	api := r.PathPrefix("/api/users").Subrouter()

	api.Use(metrics.MetricsMiddleware)
	api.Use(utils.LimitMiddleware)
	
	api.HandleFunc("", userHandler.GetUsers).Methods("GET")
	api.HandleFunc("", userHandler.CreateUser).Methods("POST")
	api.HandleFunc("/{id}", userHandler.GetUserByID).Methods("GET")
	api.HandleFunc("/{id}", userHandler.UpdateUser).Methods("PUT")
	api.HandleFunc("/{id}", userHandler.DeleteUser).Methods("DELETE")
	api.HandleFunc("/upload", integrationHandler.UploadFile).Methods("POST")

	// 4. Запуск сервера
	port := ":8080"
	fmt.Printf("Server starting on port %s...\n", port)
	fmt.Printf("MinIO connected to: %s\n", minioEndpoint) 

	if err := http.ListenAndServe(port, r); err != nil {
		log.Fatal(err)
	}
}