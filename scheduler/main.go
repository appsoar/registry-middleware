package main

import (
	"net/http"
	//"scheduler"
	//	"scheduler/middleware"
	"scheduler/log"
	"scheduler/routes"
)

func main() {
	router := routes.NewRouter()

	//	n := middleware.New()
	//	n.UseHandler(router)
	//	n.Run(":9090")
	log.Logger.Info("scheduler starts")
	http.ListenAndServe(":9090", router)

}
