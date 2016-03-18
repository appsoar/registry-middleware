package routes

import (
	"github.com/gorilla/mux"
	"golang.org/x/net/websocket"
	"net/http"
	"scheduler/handler"
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
	router.NotFoundHandler = http.HandlerFunc(handler.NotFoundHandler)

	routes = append(routes, nsRoutes...)
	routes = append(routes, repoRoutes...)
	routes = append(routes, accountRoutes...)

	for _, route := range routes {
		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(route.HandlerFunc)
	}

	router.Path("/api/v0/sysinfo").Handler(websocket.Handler(handler.GetSysInfo))
	router.Path("/api/v0/logs").Handler(websocket.Handler(handler.GetLog))
	router.Path("/api/v0/stats").Handler(websocket.Handler(handler.GetUserStats))
	return router
}

var routes = Routes{
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
		"/api/v0/login",
		handler.LoginHandler,
	},
	Route{
		"Logout",
		"POST",
		"/api/v0/logout",
		handler.LogoutHandler,
	},
}
