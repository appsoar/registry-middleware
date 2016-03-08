package scheduler

import (
	"github.com/gorilla/mux"
	"net/http"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

func NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	router.NotFoundHandler = http.HandlerFunc(NotFound)
	for _, route := range routes {
		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(route.HandlerFunc)
	}

	return router
}

var routes = Routes{
	/*
		Route{
			"List",
			"GET",
			"/images",
			ListImages,
		},
	*/
	/*
		Route{
			"Login",
			"GET",
			"/v2/login",
			ShowLogin,
		},*/
	Route{
		"Login",
		"POST",
		"/v2/login",
		Login,
	},
}