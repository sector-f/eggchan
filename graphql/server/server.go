package main

import (
	"database/sql"
	"fmt"
	log "log"
	http "net/http"
	os "os"

	handler "github.com/99designs/gqlgen/handler"
	_ "github.com/lib/pq"
	graphql "github.com/sector-f/eggchan/graphql"
	"github.com/sector-f/eggchan/postgres"
)

const defaultPort = "8080"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	connectionString := fmt.Sprintf("host=127.0.0.1 dbname=eggchan sslmode=disable")
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		log.Fatal(err)
	}

	service := postgres.EggchanService{db}
	http.Handle("/", handler.Playground("GraphQL playground", "/query"))
	http.Handle("/query", handler.GraphQL(graphql.NewExecutableSchema(graphql.Config{Resolvers: &graphql.Resolver{&service}})))

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
