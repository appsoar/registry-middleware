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

	Route{
		"Namespace",
		"PUT",
		"/api/v0/namespace",
		handler.UpdateNs,
	},

	Route{
		"Namespace",
		"DELETE",
		"/api/v0/namespace/{namespace_name}",
		handler.DeleteNs,
	},

	//UserGroup
	Route{
		"Namespace",
		"GET",
		"/api/v0/grps/{namespace}",
		handler.GetNsUgroup,
	},
	Route{
		"Namespace",
		"POST",
		"/api/v0/grp",
		handler.AddUgroup,
	},

	Route{
		"Namespace",
		"GET",
		"/api/v0/grp/{group_id}",
		handler.GetUgroup,
	},

	Route{
		"Namespace",
		"PUT",
		"/api/v0/grp",
		handler.UpdateUgroup,
	},

	Route{
		"Namespace",
		"DELETE",
		"/api/v0/grp/{group_id}",
		handler.DeleteUgroup,
	},
}
