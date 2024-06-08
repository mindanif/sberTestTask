package main

import (
	"database/sql"
	"github.com/go-chi/chi/v5"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	_ "sberTestTask/docs"
	"sberTestTask/internal/config"
	"sberTestTask/internal/todo/delivery/api"
	"sberTestTask/internal/todo/repository/postgres"
	"sberTestTask/internal/todo/service"
)

// @title Swagger Example API
// @version 1.0
// @description This is a sample server Petstore server.
// @termsOfService http://swagger.io/terms/

// @host localhost:8080
// @BasePath /
func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	db, err := sql.Open("postgres", cfg.Database.URL)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	repo := postgres.NewPostgresRepository(db)
	uc := service.NewTodoUsecase(repo)
	handler := api.NewHandler(uc)
	r := chi.NewRouter()

	api.RegisterRoutes(r, handler)

	log.Fatal(http.ListenAndServe(":"+cfg.Server.Port, r))
}
