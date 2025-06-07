package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/emday4prez/fs-go/internal/api"
)

func main() {
	router := api.NewRouter()

	port := ":8080"
	fmt.Printf("Server is starting on port %s\n", port)

	err := http.ListenAndServe(port, router)
	if err != nil {
		log.Fatal("Failed to start server: ", err)
	}
}
