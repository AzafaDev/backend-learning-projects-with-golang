package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"todo-list-api/internal/config"
	"todo-list-api/internal/database"
	"todo-list-api/internal/repository"
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
	repository.NewUserRepository(db)
	fmt.Println("server is running on port:", cfg.Port)
	http.ListenAndServe(":"+cfg.Port, mux)

}
