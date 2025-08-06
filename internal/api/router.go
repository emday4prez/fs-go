package api

import (
	"log/slog"
	"net/http"

	"github.com/emday4prez/fs-go/internal/auth"
	"github.com/emday4prez/fs-go/internal/config"
	"github.com/emday4prez/fs-go/internal/file"
	"github.com/emday4prez/fs-go/internal/user"
)

func NewRouter(fs *file.FileService, us *user.Service, as *auth.Service, cfg *config.Config, log *slog.Logger) http.Handler {
	server := NewServer(fs, us, as, cfg, log)
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
	mux.HandleFunc("/api/files", server.ListFilesAPI)
	mux.HandleFunc("/api/register", server.RegisterUserHandler)
	mux.HandleFunc("/api/login", server.LoginHandler)

	return mux
}
