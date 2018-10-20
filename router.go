package main

import (
	"github.com/gorilla/mux"
	"net/http"
)

func NewRouter() *mux.Router {
	router := mux.NewRouter()
	for _, route := range routes {
		var handler http.Handler
		handler = route.HandlerFunc
		handler = Logger(handler)

		router.
			Methods(route.Method).
			Path(route.Pattern).
			Handler(handler)

		router.
			Methods(route.Method).
			Path(route.Pattern + "/").
			Handler(handler)
	}

	return router
}
