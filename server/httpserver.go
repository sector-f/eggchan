package server

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/sector-f/eggchan"
)

type HttpServer struct {
	Router       *mux.Router
	BoardService eggchan.BoardService
	AuthService  eggchan.AuthService
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

		Route{"POST", "/new/boards/{board}", s.createBoard, true, "create_board"},
	}

	router := mux.NewRouter().StrictSlash(true)
	router.NotFoundHandler = http.HandlerFunc(handleNotFound)

	for _, route := range routes {
		var handler http.Handler
		handler = route.HandlerFunc

		if route.AuthRequired {
			handler = s.auth(handler, route.Permission)
		}

		handler = handlers.LoggingHandler(os.Stdout, handlers.ProxyHeaders(handler))
		router.Methods(route.Method).Path(route.Pattern).Handler(handler)

		if route.Method == "GET" {
			router.Methods("HEAD").Path(route.Pattern).Handler(handler)
		}
	}

	s.Router = router
}

func (s *HttpServer) Run(addr string) {
	log.Fatal(http.ListenAndServe(addr, s.Router))
}
