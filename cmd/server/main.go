package main

import (
	"log"
	"net/http"

	"github.com/emday4prez/fs-go/internal/api"
	"github.com/emday4prez/fs-go/internal/config"
	"github.com/emday4prez/fs-go/internal/service"
	"github.com/emday4prez/fs-go/internal/storage"
)

func main() {
	cfg := config.Load()

	localStorage, err := storage.NewLocalStorage(cfg.UploadDir)
	if err != nil {
		log.Fatalf("failed to create local storage: %v", err)
	}
	fileService := service.NewFileService(localStorage)

	router := api.NewRouter(fileService, cfg)

	log.Printf("Server is starting on port %s", cfg.Port)
	if err := http.ListenAndServe(cfg.Port, router); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
