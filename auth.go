package main

import (
	"net/http"
)

func (a *Server) auth(inner http.Handler, perm string) http.Handler {
	return a.checkPassword(a.checkPermissions(inner, perm))
}

func (a *Server) checkPassword(inner http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username, password, ok := r.BasicAuth()
		if !ok {
			respondWithError(w, http.StatusUnauthorized, "Authentication required")
			return
		}

		pw_is_correct, err := getUserAuthentication(a.DB, username, []byte(password))
		if err != nil {
			respondWithError(w, http.StatusUnauthorized, "Bad username or password")
			return
		}

		if !pw_is_correct {
			respondWithError(w, http.StatusUnauthorized, "Bad username or password")
		} else {
			inner.ServeHTTP(w, r)
		}
	})
}

func (a *Server) checkPermissions(inner http.Handler, perm string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username, _, ok := r.BasicAuth()
		if !ok {
			respondWithError(w, http.StatusUnauthorized, "Authentication required")
			return
		}

		has_perm, err := getUserAuthorization(a.DB, username, perm)
		if err != nil || has_perm == false {
			respondWithError(w, http.StatusUnauthorized, "Insufficient permissions")
			return
		}

		inner.ServeHTTP(w, r)
	})
}
