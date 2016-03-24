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
		"GET",
		"/api/v0/account/{account}",
		handler.GetUserAccount,
	},

	Route{
		"Account",
		"POST",
		"/api/v0/account",
		handler.AddAccount,
	},

	Route{
		"Account",
		"PUT",
		"/api/v0/account",
		handler.UpdateAccount,
	},

	Route{
		"Account",
		"DELETE",
		"/api/v0/account/{user_id}",
		handler.DeleteAccount,
	},
}
