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

func New(e *eggchan.EggchanService) HttpServer {
	var routes = Routes{
		Route{"GET", "/categories", e.getCategories, false, ""},
		Route{"GET", "/categories/{category}", e.showCategory, false, ""},

		Route{"GET", "/boards", e.getBoards, false, ""},
		Route{"GET", "/boards/{board}", e.showBoard, false, ""},
		Route{"POST", "/boards/{board}", e.postThread, false, ""},

		Route{"GET", "/boards/{board}/{thread}", e.showThread, false, ""},
		Route{"POST", "/boards/{board}/{thread}", e.postReply, false, ""},

		Route{"DELETE", "/boards/{board}/threads/{thread}", e.deleteThread, true, "delete_thread"},
		Route{"DELETE", "/boards/{board}/comments/{comment}", e.deleteComment, true, "delete_post"},
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
