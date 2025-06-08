package api

import (
	"fmt"
	"html/template"
	"log"
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
