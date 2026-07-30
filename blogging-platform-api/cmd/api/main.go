package main

import (
	"blogging-platform-api/internal/config"
	"blogging-platform-api/internal/database"
	"blogging-platform-api/internal/post"
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	ctx := context.Background()

	pool, err := database.Connect(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("failed to connect db: %v", err)
	}
	defer pool.Close()

	repo := post.NewRepository(pool)
	service := post.NewService(repo)
	handler := post.NewHandler(service)

	mux := http.NewServeMux()
	mux.HandleFunc("POST /posts", handler.CreatePost)
	mux.HandleFunc("GET /posts/{id}", handler.GetByID)
	mux.HandleFunc("PUT /posts/{id}", handler.UpdatePost)
	mux.HandleFunc("DELETE /posts/{id}", handler.DeletePost)
	mux.HandleFunc("GET /posts", handler.SearchPost)

	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: mux,
	}

	go func() {
		log.Printf("server running on: %s", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	log.Println("shutting down server...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("server forced to shutdown: %v", err)
	}

	log.Println("server exited")
}
