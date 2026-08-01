package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"todo-list-api/internal/config"
	"todo-list-api/internal/database"
	"todo-list-api/internal/handlers"
	"todo-list-api/internal/middleware"
	"todo-list-api/internal/repository"
	"todo-list-api/internal/services"

	"github.com/rs/zerolog/log"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load config")
	}

	db, err := database.ConnectDB(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect to database")
	}
	defer db.Close()

	userRepo := repository.NewUserRepository(db)
	userService := services.NewUserService(userRepo, cfg)
	userHandler := handlers.NewUserHandler(userService)

	todoRepo := repository.NewTodoRepository(db)
	todoService := services.NewTodoService(todoRepo)
	todoHandler := handlers.NewTodoHandler(todoService)

	mux := http.NewServeMux()
	mux.HandleFunc("POST /register", userHandler.Register)
	mux.HandleFunc("POST /login", userHandler.Login)

	mux.Handle("POST /todos", middleware.AuthMiddleware(http.HandlerFunc(todoHandler.CreateTodos), *cfg))
	mux.Handle("PUT /todos/{id}", middleware.AuthMiddleware(http.HandlerFunc(todoHandler.UpdateTodo), *cfg))
	mux.Handle("DELETE /todos/{id}", middleware.AuthMiddleware(http.HandlerFunc(todoHandler.DeleteTodo), *cfg))
	mux.Handle("GET /todos", middleware.AuthMiddleware(http.HandlerFunc(todoHandler.GetTodos), *cfg))
	mux.Handle("GET /todos/{id}", middleware.AuthMiddleware(http.HandlerFunc(todoHandler.GetTodoByID), *cfg))

	srv := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	go func() {
		log.Info().Str("port", cfg.Port).Msg("server is running")
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal().Err(err).Msg("server failed")
		}
	}()

	<-ctx.Done()
	log.Info().Msg("shutting down server")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Error().Err(err).Msg("server shutdown failed")
	}
}
