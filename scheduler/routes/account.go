package routes

import (
	"scheduler/handler"
)

var accountRoutes = Routes{
	Route{
		"Account",
		"GET",
		"/api/v0/accounts",
		handler.GetAccounts,
	},

	Route{
		"Account",
		"POST",
		"/api/v0/account",
		handler.GetAccounts,
	},
}
