package main

import (
	"database/sql"
	"log"

	"github.com/ericolvr/sec-backend/config"
	"github.com/ericolvr/sec-backend/internal/core/services"
	"github.com/ericolvr/sec-backend/internal/infrastructure/database"
	httpInfra "github.com/ericolvr/sec-backend/internal/infrastructure/http"
	"github.com/ericolvr/sec-backend/internal/infrastructure/http/routes"
	"github.com/ericolvr/sec-backend/internal/interfaces/api"
	_ "github.com/lib/pq"
)

func main() {
	cfg := config.Load()

	db, err := sql.Open("postgres", cfg.GetDatabaseURL())
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("Error pinging database: %v", err)
	}

	log.Println("Database connected")

	// Repository
	userRepo := database.NewUserRepository(db)

	// Service
	userService := services.NewUserService(userRepo)

	// Handler
	userHandler := api.NewUserHandler(userService)

	// Routes
	userRoutes := routes.NewUserRoutes(userHandler)

	// Router
	router := httpInfra.NewRouter(userRoutes)

	// Server
	server := httpInfra.NewServer(router, cfg.Server.Port)

	log.Printf("Server started on port %s", cfg.Server.Port)
	if err := server.Start(); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
