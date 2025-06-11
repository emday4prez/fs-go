package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/emday4prez/fs-go/internal/api"
	"github.com/emday4prez/fs-go/internal/config"
	"github.com/emday4prez/fs-go/internal/service"
)

func main() {
	cfg := config.Load()

	fileService := service.NewFileService(cfg.UploadDir)

	router := api.NewRouter(fileService, cfg)

	fmt.Printf("Server is starting on port %s\n", cfg.Port)

	err := http.ListenAndServe(cfg.Port, router)
	if err != nil {
		log.Fatal("Failed to start server: ", err)
	}
}
