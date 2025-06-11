package api

import (
	"net/http"

	"github.com/emday4prez/fs-go/internal/config"
	"github.com/emday4prez/fs-go/internal/service"
)

func NewRouter(fs *service.FileService, cfg *config.Config) http.Handler {
	server := NewServer(fs, cfg)
	mux := http.NewServeMux()

	mux.HandleFunc("/", server.WelcomeHandler)

	uploadMultiplexer := func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			server.ShowUploadPage(w, r)
		} else if r.Method == http.MethodPost {
			server.UploadHandler(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}

	mux.HandleFunc("/upload", uploadMultiplexer)
	mux.HandleFunc("/list", server.ShowListPage)
	mux.HandleFunc("/download/", server.DownloadHandler)

	return mux
}
