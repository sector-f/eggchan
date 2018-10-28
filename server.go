package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"log"
	"net/http"
)

type Route struct {
	Method       string
	Pattern      string
	HandlerFunc  http.HandlerFunc
	AuthRequired bool
	Permission   string
}

type Routes []Route

type Server struct {
	Router *mux.Router
	DB     *sql.DB
}

func (a *Server) Initialize(user, password, dbname string) {
	connectionString := fmt.Sprintf("host=127.0.0.1 dbname=%s sslmode=disable", dbname)

	var err error
	a.DB, err = sql.Open("postgres", connectionString)
	if err != nil {
		log.Fatal(err)
	}

	var routes = Routes{
		Route{
			"GET",
			"/categories",
			a.getCategories,
			false,
			"",
		},
		Route{
			"GET",
			"/categories/{category}",
			a.showCategory,
			false,
			"",
		},
		Route{
			"GET",
			"/boards",
			a.getBoards,
			false,
			"",
		},
		Route{
			"GET",
			"/boards/{board}",
			a.showBoard,
			false,
			"",
		},
		Route{
			"POST",
			"/boards/{board}",
			a.postThread,
			false,
			"",
		},
		Route{
			"POST",
			"/boards/{board}/{thread}",
			a.postReply,
			false,
			"",
		},
		Route{
			"GET",
			"/boards/{board}/{thread}",
			a.showThread,
			false,
			"",
		},
	}

	router := mux.NewRouter()
	router.NotFoundHandler = http.HandlerFunc(handleNotFound)

	for _, route := range routes {
		var handler http.Handler
		handler = route.HandlerFunc

		if route.AuthRequired {
			handler = a.auth(handler, route.Permission)
		}

		handler = Logger(handler)
		router.Methods(route.Method).Path(route.Pattern).Handler(handler)
		router.Methods(route.Method).Path(route.Pattern + "/").Handler(handler)
	}

	a.Router = router
}

func (a *Server) Run(addr string) {
	log.Fatal(http.ListenAndServe(addr, a.Router))
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
