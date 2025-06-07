package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/emday4prez/fs-go/internal/api"
	"github.com/emday4prez/fs-go/internal/service"
)

func main() {
	fileService := service.NewFileService()

	router := api.NewRouter(fileService)

	port := ":8080"
	fmt.Printf("Server is starting on port %s\n", port)

	err := http.ListenAndServe(port, router)
	if err != nil {
		log.Fatal("Failed to start server: ", err)
	}
}
