package scheduler

import (
	"github.com/gorilla/mux"
	"golang.org/x/net/websocket"
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

	router.Path("/v2/sysinfo").Handler(websocket.Handler(GetSysInfo))
	router.Path("/v2/logs").Handler(websocket.Handler(GetLog))
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
	/*websocket can't use
	Route{
		"Login",
		"GET",
		"/v2/sysinfo",
		websocket.Handler(GetSysInfo),
	},*/
	Route{
		"Login",
		"POST",
		"/v2/login",
		Login,
	},
	Route{
		"Logout",
		"GET",
		"/v2/logout",
		Logout,
	},
	Route{
		"Test",
		"GET",
		"/v2/test",
		Test,
	},
}
