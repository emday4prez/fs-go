package api

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/emday4prez/fs-go/internal/config"
	"github.com/emday4prez/fs-go/internal/service"
)

type Server struct {
	fileService *service.FileService
	config      *config.Config
}

func NewServer(fs *service.FileService, cfg *config.Config) *Server {
	return &Server{
		fileService: fs,
		config:      cfg,
	}
}

func (s *Server) WelcomeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the go file server.")
}

func (s *Server) ShowUploadPage(w http.ResponseWriter, r *http.Request) {
	tmp, err := template.ParseFiles("web/template/upload.html")
	if err != nil {
		log.Printf("Error parsing template: %v", err)
		http.Error(w, "There was a problem preparing the page.", http.StatusInternalServerError)
		return
	}
	err = tmp.Execute(w, nil)
	if err != nil {
		log.Printf("error executing template: %v", err)
		http.Error(w, "there was a problem rendering the page", http.StatusInternalServerError)
	}
}

func (s *Server) UploadHandler(w http.ResponseWriter, r *http.Request) {

	//parse headers, max size 10mb
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, "file too big, max 10MB", http.StatusBadRequest)
		return
	}

	// get file from 'name' html attribute
	_, fileHeader, err := r.FormFile("file")
	if err != nil {
		log.Printf("Error retrieving file from: %v", err)
		http.Error(w, "error retrieving the file", http.StatusInternalServerError)
		return
	}

	// call service
	err = s.fileService.SaveFile(fileHeader)
	if err != nil {
		log.Printf("error saving file: %v", err)
		http.Error(w, "error saving the file", http.StatusInternalServerError)
		return
	}

	// send response
	fmt.Fprintf(w, "file '%s' uploaded successfully!", fileHeader.Filename)

}

func (s *Server) ShowListPage(w http.ResponseWriter, r *http.Request) {
	filenames, err := s.fileService.ListFiles()
	if err != nil {
		log.Printf("Error listing files: %v", err)
		http.Error(w, "Could not retrieve file list.", http.StatusInternalServerError)
		return
	}

	tmp, err := template.ParseFiles("web/template/list.html")
	if err != nil {
		log.Printf("Error parsing template: %v", err)
		http.Error(w, "There was a problem preparing the page.", http.StatusInternalServerError)
		return
	}
	err = tmp.Execute(w, filenames)
	if err != nil {
		log.Printf("error executing template: %v", err)
		http.Error(w, "there was a problem rendering the page", http.StatusInternalServerError)
	}
}

func (s *Server) DownloadHandler(w http.ResponseWriter, r *http.Request) {
	filename := strings.TrimPrefix(r.URL.Path, "/download/")
	if filename == "" {
		http.NotFound(w, r)
		return
	}

	filePath := filepath.Join(s.config.UploadDir, filename)

	http.ServeFile(w, r, filePath)
}
