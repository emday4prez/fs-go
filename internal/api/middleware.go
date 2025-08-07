package api

import (
	"net/http"
	"strings"
)

func (s *Server) authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//get auth header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			s.logger.Warn("Missing Authorization header")
			http.Error(w, "Authorization header is required", http.StatusUnauthorized)
			return
		}

		// Bearer <token>
		headerParts := strings.Split(authHeader, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			http.Error(w, "Invalid Authorization header format", http.StatusUnauthorized)
			return
		}
		tokenString := headerParts[1]

		_, err := s.authService.ValidateJWT(tokenString)
		if err != nil {
			s.logger.Error("Invalid token", "error", err)
			http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
