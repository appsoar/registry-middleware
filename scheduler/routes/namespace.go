package routes

import (
	"scheduler/handler"
)

var nsRoutes = Routes{
	Route{
		"Namespace",
		"GET",
		"/api/v0/namespaces",
		handler.GetAllNsHandler,
	},
	Route{
		"Namespace",
		"GET",
		"/api/v0/namespace/{namespace}",
		handler.GetSpecNsHandler,
	},

	//UserGroup
	Route{
		"Namespace",
		"GET",
		"/api/v0/grp/{namespace}",
		handler.GetNsUgroup,
	},
	Route{
		"Namespace",
		"POST",
		"/api/v0/grp",
		handler.AddUgroup,
	},
}
