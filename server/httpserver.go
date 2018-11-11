package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sector-f/eggchan"
)

type HttpServer struct {
	Router         *mux.Router
	EggchanService *eggchan.EggchanService
}

type Route struct {
	Method       string
	Pattern      string
	HandlerFunc  http.HandlerFunc
	AuthRequired bool
	Permission   string
}

type Routes []Route

func (s *HttpServer) Initialize() {
	var routes = Routes{
		Route{"GET", "/categories", s.getCategories, false, ""},
		Route{"GET", "/categories/{category}", s.showCategory, false, ""},

		Route{"GET", "/boards", s.getBoards, false, ""},
		Route{"GET", "/boards/{board}", s.showBoard, false, ""},
		Route{"POST", "/boards/{board}", s.postThread, false, ""},

		Route{"GET", "/boards/{board}/{thread}", s.showThread, false, ""},
		Route{"POST", "/boards/{board}/{thread}", s.postReply, false, ""},

		Route{"DELETE", "/boards/{board}/threads/{thread}", s.deleteThread, true, "delete_thread"},
		Route{"DELETE", "/boards/{board}/comments/{comment}", s.deleteComment, true, "delete_post"},
	}

	router := mux.NewRouter().StrictSlash(true)
	router.NotFoundHandler = http.HandlerFunc(handleNotFound)

	for _, route := range routes {
		var handler http.Handler
		handler = route.HandlerFunc

		if route.AuthRequired {
			handler = s.auth(handler, route.Permission)
		}

		handler = Logger(handler)
		router.Methods(route.Method).Path(route.Pattern).Handler(handler)
	}

	s.Router = router
}
