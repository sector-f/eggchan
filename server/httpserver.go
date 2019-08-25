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
	AdminService eggchan.AdminService
	AuthService  eggchan.AuthService
	routes       Routes
}

type Route struct {
	Method       string
	Pattern      string
	Description  string
	HandlerFunc  http.HandlerFunc
	AuthRequired bool
	Permission   string
}

type Routes []Route

func (s *HttpServer) Initialize() {
	var routes = Routes{
		Route{"GET", "/", "List routes", s.index, false, ""},

		Route{"GET", "/categories", "List categories", s.getCategories, false, ""},
		Route{"GET", "/categories/{category}", "List boards in a specific category", s.showCategory, false, ""},

		Route{"GET", "/boards", "List boards", s.getBoards, false, ""},
		Route{"GET", "/boards/{board}", "List threads in a specific board", s.showBoard, false, ""},
		Route{"POST", "/boards/{board}", "Post to a specific board", s.postThread, false, ""},

		Route{"GET", "/boards/{board}/{thread}", "Show a specific thread", s.showThread, false, ""},
		Route{"POST", "/boards/{board}/{thread}", "Post to a specific thread", s.postReply, false, ""},

		Route{"DELETE", "/boards/{board}/threads/{thread}", "Delete a specific thread", s.deleteThread, true, "delete_thread"},
		Route{"DELETE", "/boards/{board}/comments/{comment}", "Delete a specific comment", s.deleteComment, true, "delete_post"},

		Route{"POST", "/new/boards/{board}", "Create a new board", s.createBoard, true, "create_board"},
	}

	s.routes = routes

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

func prettyHandler(h http.Handler, pretty bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.ServeHTTP(w, r)
	}
}

func (s *HttpServer) Run(addr string) {
	log.Fatal(http.ListenAndServe(addr, s.Router))
}
