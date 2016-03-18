package routes

import (
	"scheduler/handler"
)

var nsRoutes = Routes{
	Route{
		"Namespace",
		"GET",
		"/api/v0/namespaces",
		handler.NamespacesGetHandler,
	},
	Route{
		"Namespace",
		"GET",
		"/api/v0/namespace/{namespace}",
		handler.NamespaceSpecificGetHandler,
	},
}
