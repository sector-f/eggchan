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

		pw_is_correct, err := s.AuthService.ValidatePassword(username, []byte(password))
		if !pw_is_correct || err != nil {
			respondWithError(w, http.StatusUnauthorized, "Bad username or password")
			return
		}

		has_perm, err := s.AuthService.CheckPermission(username, permission)
		if !has_perm || err != nil {
			respondWithError(w, http.StatusUnauthorized, "Insufficient permissions")
			return
		}

		inner.ServeHTTP(w, r)
	})
}
