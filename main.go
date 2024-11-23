package main

import (
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/devtiagoneves/api-restful/api"
	"github.com/devtiagoneves/api-restful/pkg"
)

func main() {
	if err := run(); err != nil {
		slog.Error("failed to execute code", "error", err)
		os.Exit(1)
	}

	slog.Info("all systems offline")
}

func run() error {
	db := pkg.NewDB()
	handler := api.NewHandler(db)

	s := http.Server{
		ReadTimeout:  10 * time.Second,
		WriteTimeout: time.Minute,
		IdleTimeout:  10 * time.Second,
		Addr:         ":8080",
		Handler:      handler,
	}

	if err := s.ListenAndServe(); err != nil {
		return err
	}

	return nil
}
