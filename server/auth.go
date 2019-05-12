package server

import (
	"net/http"
)

func (s *HttpServer) auth(inner http.Handler, permission string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username, password, ok := r.BasicAuth()
		if !ok {
			respondWithError(w, http.StatusUnauthorized, "Authentication required")
			return
		}

		hasPermission, err := s.AuthService.CheckAuth(username, []byte(password), permission)
		if !hasPermission || err != nil {
			respondWithError(w, http.StatusUnauthorized, "Permission denied")
			return
		}

		inner.ServeHTTP(w, r)
	})
}
