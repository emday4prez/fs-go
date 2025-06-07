package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", welcomeHandler)

	port := ":8080"
	fmt.Printf("Server is starting on port %s\n", port)

	err := http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatal("Failed to start server: ", err)
	}
}

func welcomeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, World! Wecome to the Go File Server")
}
