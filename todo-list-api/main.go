package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"todo-list-api/internal/config"
	"todo-list-api/internal/database"
)

func main() {
	ctx := context.Background()
	mux := http.NewServeMux()
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatal("error:", err)
	}
	_, err = database.ConnectDB(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatal("error:", err)
	}
	fmt.Println("server is running on port:", cfg.Port)
	http.ListenAndServe(":"+cfg.Port, mux)

}
