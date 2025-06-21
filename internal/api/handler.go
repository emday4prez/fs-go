package api

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"strings"

	"github.com/emday4prez/fs-go/internal/config"
	"github.com/emday4prez/fs-go/internal/service"
)

type Server struct {
	fileService *service.FileService
	config      *config.Config
	logger      *slog.Logger
}

func NewServer(fs *service.FileService, cfg *config.Config, log *slog.Logger) *Server {
	return &Server{
		fileService: fs,
		config:      cfg,
		logger:      log,
	}
}

func (s *Server) WelcomeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the go file server.")
}

func (s *Server) ShowUploadPage(w http.ResponseWriter, r *http.Request) {
	tmp, err := template.ParseFiles("web/template/upload.html")
	if err != nil {
		s.logger.Error("Could not parse template", "path", "web/template/list.html", "error", err)
		http.Error(w, "There was a problem preparing the page.", http.StatusInternalServerError)
		return
	}
	err = tmp.Execute(w, nil)
	if err != nil {
		s.logger.Error("Could not execute template", "error", err)
		http.Error(w, "there was a problem rendering the page", http.StatusInternalServerError)
	}
}

func (s *Server) UploadHandler(w http.ResponseWriter, r *http.Request) {

	//parse headers, max size 10mb
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		s.logger.Error("file too big to upload", "error", err)
		http.Error(w, "file too big, max 10MB", http.StatusBadRequest)
		return
	}

	// get file from 'name' html attribute
	_, fileHeader, err := r.FormFile("file")
	if err != nil {
		s.logger.Error("error retrieving file", "error", err)
		http.Error(w, "error retrieving the file", http.StatusInternalServerError)
		return
	}

	// call service
	err = s.fileService.SaveFile(fileHeader)
	if err != nil {
		s.logger.Error("error saving file", "error", err)
		http.Error(w, "error saving the file", http.StatusInternalServerError)
		return
	}
	// send response
	fmt.Fprintf(w, "file '%s' uploaded successfully!", fileHeader.Filename)
}

func (s *Server) ShowListPage(w http.ResponseWriter, r *http.Request) {
	filenames, err := s.fileService.ListFiles()
	if err != nil {
		s.logger.Error("Could not retrieve file list", "error", err)
		http.Error(w, "Could not retrieve file list.", http.StatusInternalServerError)
		return
	}

	tmp, err := template.ParseFiles("web/template/list.html")
	if err != nil {
		s.logger.Error("Could not parse template", "path", "web/template/list.html", "error", err)
		http.Error(w, "There was a problem preparing the page.", http.StatusInternalServerError)
		return
	}
	err = tmp.Execute(w, filenames)
	if err != nil {
		s.logger.Error("error executing template", "error", err)
		http.Error(w, "there was a problem rendering the page", http.StatusInternalServerError)
	}
}

func (s *Server) DownloadHandler(w http.ResponseWriter, r *http.Request) {
	filename := strings.TrimPrefix(r.URL.Path, "/download/")
	if filename == "" {
		http.NotFound(w, r)
		return
	}

	_, filePath, err := s.fileService.GetFile(filename)
	if err != nil {
		s.logger.Error("error getting file", "error", err)
		http.Error(w, "File not found.", http.StatusNotFound)
		return
	}

	http.ServeFile(w, r, filePath)
}

func (s *Server) ListFilesAPI(w http.ResponseWriter, r *http.Request) {
	filenames, err := s.fileService.ListFiles()
	if err != nil {
		s.logger.Error("error listing files for api", "error", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)

		json.NewEncoder(w).Encode(map[string]string{"error": "Could not retrieve file list."})
		return
	}
	w.Header().Set("Content-Type", "application/json")

	err = json.NewEncoder(w).Encode(filenames)
	if err != nil {

		s.logger.Error("error encoding json", "error", err)
	}
}
