package api

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"strings"

	"github.com/emday4prez/fs-go/internal/config"
	"github.com/emday4prez/fs-go/internal/file"
	"github.com/emday4prez/fs-go/internal/user"
)

type Server struct {
	fileService *file.FileService
	userService *user.Service
	config      *config.Config
	logger      *slog.Logger
}

func NewServer(fs *file.FileService, us *user.Service, cfg *config.Config, log *slog.Logger) *Server {
	return &Server{
		fileService: fs,
		userService: us,
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

func (s *Server) RegisterUserHandler(w http.ResponseWriter, r *http.Request) {

	// in a more real world scenario parse json body
	if err := r.ParseForm(); err != nil {
		s.logger.Error("Failed to parse form", "error", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	username := r.FormValue("username")
	password := r.FormValue("password")

	if username == "" || password == "" {
		http.Error(w, "Username and password are required", http.StatusBadRequest)
		return
	}
	newUser, err := s.userService.Register(username, password)
	if err != nil {

		if err == user.ErrUsernameTaken {
			http.Error(w, err.Error(), http.StatusConflict) // 409
			return
		}

		s.logger.Error("Failed to register user", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated) // 201

	response := map[string]string{
		"id":       newUser.ID,
		"username": newUser.Username,
	}
	json.NewEncoder(w).Encode(response)
}

func (s *Server) LoginHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		s.logger.Error("failed to parse form", "error", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

}
