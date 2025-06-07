package api

import (
	"fmt"
	"net/http"

	"github.com/emday4prez/fs-go/internal/service"
)

type Server struct {
	fileService *service.FileService
}

func NewServer(fs *service.FileService) *Server {
	return &Server{
		fileService: fs,
	}
}

func (s *Server) WelcomeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the go file server.")
}
