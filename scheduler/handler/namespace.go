package handler

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"scheduler/errjson"
	"scheduler/log"
)

func namespacesGet(w http.ResponseWriter, r *http.Request) (err error) {
	user, err := getRequestUser(w, r)
	if err != nil {
		err = errjson.NewUnauthorizedError("user doesn't login")
		//errJsonReturn(w, r, e)
		return
	}

	log.Logger.Info(user + "get namespace info")
	nsJson, err := globalClient.GetNamespaces()
	if err != nil {
		err = errjson.NewInternalServerError("can't get ns info")
		return
	}
	fmt.Fprintf(w, string(nsJson))
	return
}

func namespaceGetSpecific(w http.ResponseWriter, r *http.Request) (err error) {
	user, err := getRequestUser(w, r)
	if err != nil {
		err = errjson.NewUnauthorizedError("user doesn't login")
		//errJsonReturn(w, r, e)
		return
	}

	vars := mux.Vars(r)
	ns := vars["namespace"]
	if len(ns) == 0 {
		err = errjson.NewNotValidEntityError("invalid namespace")
		return
	}

	log.Logger.Info(user + " get " + ns + " namespace info")

	nsJson, err := globalClient.GetSpecificNamespace(ns)
	_, err = globalClient.GetSpecificNamespace(ns)
	if err != nil {
		err = errjson.NewInternalServerError("can't get ns info")
		return
	}
	fmt.Fprintf(w, string(nsJson))
	return
}
