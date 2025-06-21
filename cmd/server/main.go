package main

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/emday4prez/fs-go/internal/api"
	"github.com/emday4prez/fs-go/internal/config"
	"github.com/emday4prez/fs-go/internal/service"
	"github.com/emday4prez/fs-go/internal/storage"
)

func main() {

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	cfg := config.Load()

	localStorage, err := storage.NewLocalStorage(cfg.UploadDir)
	if err != nil {
		logger.Error("Failed to create local storage", "error", err)
		os.Exit(1)
	}
	fileService := service.NewFileService(localStorage)

	router := api.NewRouter(fileService, cfg, logger)

	logger.Info("Server is starting", "port", cfg.Port)
	if err := http.ListenAndServe(cfg.Port, router); err != nil {
		logger.Error("Failed to start server", "error", err)
		os.Exit(1)
	}
}
