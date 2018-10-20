package main

import "net/http"

type Route struct {
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

var routes = Routes{
	Route{
		"GET",
		"/categories",
		getCategories,
	},
	Route{
		"GET",
		"/categories/{category}",
		showCategory,
	},
	Route{
		"GET",
		"/boards",
		getBoards,
	},
	Route{
		"GET",
		"/boards/{board}",
		showBoard,
	},
	Route{
		"GET",
		"/boards/{board}/{thread}",
		showThread,
	},
}
