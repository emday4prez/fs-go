package api

import (
	"net/http"

	"github.com/emday4prez/fs-go/internal/service"
)

func NewRouter(fs *service.FileService) http.Handler {
	server := NewServer(fs)

	mux := http.NewServeMux()

	mux.HandleFunc("/", server.WelcomeHandler)

	return mux
}
