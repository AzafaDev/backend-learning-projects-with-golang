package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"todo-list-api/internal/config"
	"todo-list-api/internal/database"
	"todo-list-api/internal/handlers"
	"todo-list-api/internal/repository"
	"todo-list-api/internal/services"
)

func main() {
	ctx := context.Background()
	mux := http.NewServeMux()
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatal("error:", err)
	}
	db, err := database.ConnectDB(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatal("error:", err)
	}
	userRepo := repository.NewUserRepository(db)

	userService := services.NewUserService(userRepo, cfg)

	userHandler := handlers.NewUserHandler(userService)

	mux.HandleFunc("POST /register", func(w http.ResponseWriter, r *http.Request) {
		userHandler.Register(w, r, ctx)
	})

	mux.HandleFunc("POST /login", func(w http.ResponseWriter, r *http.Request) {
		userHandler.Login(w, r, ctx)
	})
	fmt.Println("server is running on port:", cfg.Port)
	http.ListenAndServe(":"+cfg.Port, mux)

}
