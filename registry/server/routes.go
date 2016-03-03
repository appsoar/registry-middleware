package server

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
	Route{
		"List",
		"GET",
		"/images",
		ListImages,
	},

	Route{
		"Image",
		"GET",
		//"/image/{image:[a-Z0-9]}/tag",
		//注意:这会出现-在句首的错误
		"/{image:[a-zA-Z0-9]+[-a-zA-Z0-9]*(/[a-zA-Z0-9]+[-a-zA-Z0-9]*)*}/tag",
		ListImageTags,
	},
	Route{
		"Image",
		"Delete",
		"/{image:[a-zA-Z0-9]+[-a-zA-Z0-9]*(/[a-zA-Z0-9]+[-a-zA-Z0-9]*)*}/{tag}",
		DeleteImageTag,
	},
}
