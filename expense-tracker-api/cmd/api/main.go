package main

import (
	"context"
	"expense-tracker-api/internal/config"
	"expense-tracker-api/internal/database"
	"expense-tracker-api/internal/user"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rs/zerolog/log"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	mux := http.NewServeMux()

	cfg, err := config.LoadEnv()
	if err != nil {
		log.Fatal().Err(err).Msg("error in loadEnv")
	}

	db, err := database.NewDatabase(cfg.DatabaseURL, ctx)

	userRepo := user.NewUserRepository(db)
	userService := user.NewUserService(userRepo)
	userHandler := user.NewUserHandler(userService, cfg)

	mux.HandleFunc("POST /register", userHandler.Register)
	mux.HandleFunc("POST /login", userHandler.Login)

	mux.HandleFunc("GET /slow", func(w http.ResponseWriter, r *http.Request) {
		log.Info().Msg("request started")

		time.Sleep(10 * time.Second)

		log.Info().Msg("request finished")

		w.Write([]byte("done"))
	})

	server := http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	go func() {
		log.Info().Str("port", cfg.Port).Msg("server is running")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("server failed to listen and serve")
		}
	}()

	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Error().Err(err).Msg("server is failed to shutdown")
	}
}
