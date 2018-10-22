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
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

type Server struct {
	Router    *mux.Router
	DB        *sql.DB
	BumpLimit int
}

func (a *Server) Initialize(user, password, dbname string) {
	a.BumpLimit = 300

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
		},
		Route{
			"GET",
			"/categories/{category}",
			a.showCategory,
		},
		Route{
			"GET",
			"/boards",
			a.getBoards,
		},
		Route{
			"GET",
			"/boards/{board}",
			a.showBoard,
		},
		Route{
			"POST",
			"/boards/{board}",
			a.postThread,
		},
		Route{
			"POST",
			"/boards/{board}/{thread}",
			a.postReply,
		},
		Route{
			"GET",
			"/boards/{board}/{thread}",
			a.showThread,
		},
	}

	router := mux.NewRouter()
	router.NotFoundHandler = http.HandlerFunc(handleNotFound)

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
