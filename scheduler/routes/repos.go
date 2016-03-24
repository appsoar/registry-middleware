package routes

import (
	"scheduler/handler"
)

var repoRoutes = Routes{
	Route{
		"Repos",
		"GET",
		"/api/v0/repositories",
		handler.GetReposHandler,
	},
	Route{
		"Repos",
		"GET",
		"/api/v0/repository/{repoName:[-a-zA-Z0-9]+(/[-a-zA-Z0-9]+)*}",
		handler.ListRepoTagsHandler,
	},
	Route{
		"Repos",
		"GET",
		"/api/v0/repositories/{namespace}",
		handler.GetNsReposHandler,
	},
	Route{
		"Repos",
		"GET",
		"/api/v0/repositories/user/{user}",
		handler.GetUserReposHandler,
	},
	Route{
		"Repos",
		"GET",
		"/api/v0/tag/{repoName:[-a-zA-Z0-9]+(/[-a-zA-Z0-9]+)*}/{tagName}",
		handler.GetTagImageHandler,
	},
}
