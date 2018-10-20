package main

import "net/http"

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

var routes = Routes{
	Route{
		"List Categories",
		"GET",
		"/categories",
		getCategories,
	},
	Route{
		"Show Category",
		"GET",
		"/categories/{category}",
		showCategory,
	},
	Route{
		"List Boards",
		"GET",
		"/boards",
		getBoards,
	},
	Route{
		"Show Board",
		"GET",
		"/boards/{board}",
		showBoard,
	},
	Route{
		"Show Thread",
		"GET",
		"/boards/{board}/{thread}",
		showThread,
	},
}
