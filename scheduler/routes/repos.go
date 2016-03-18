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
		"/api/v0/repository/{usernameOrNamespace}/{repoName}",
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
		"/api/v0/repository/{usernameOrNamespace}/{repoName}/{tagName}",
		handler.GetTagImageHandler,
	},
}
