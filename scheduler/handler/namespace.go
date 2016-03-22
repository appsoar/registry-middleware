package handler

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"scheduler/client/database"
	"scheduler/errjson"
	"scheduler/log"
)

func GetAllNs(w http.ResponseWriter, r *http.Request) (err error) {
	user, err := getRequestUser(w, r)
	if err != nil {
		err = errjson.NewUnauthorizedError("user doesn't login")
		//errJsonReturn(w, r, e)
		return
	}

	log.Logger.Info(user + "get namespace info")
	nsJson, err := globalClient.GetNamespaces()
	if err != nil {
		err = checkDbErr(err)
		return
	}
	fmt.Fprintf(w, string(nsJson))
	return
}

func getSpecNs(w http.ResponseWriter, r *http.Request) (err error) {
	user, err := getRequestUser(w, r)
	if err != nil {
		err = errjson.NewUnauthorizedError("user doesn't login")
		//errJsonReturn(w, r, e)
		return
	}

	vars := mux.Vars(r)
	ns := vars["namespace"]
	if len(ns) == 0 {
		err = errjson.NewErrForbidden("invalid namespace")
		return
	}

	log.Logger.Info(user + " get " + ns + " namespace info")

	nsJson, err := globalClient.GetSpecificNamespace(ns)
	if err != nil {
		err = checkDbErr(err)
		return
	}
	fmt.Fprintf(w, string(nsJson))
	return
}

func getNsUgroup(w http.ResponseWriter, r *http.Request) (err error) {
	user, err := getRequestUser(w, r)
	if err != nil {
		err = errjson.NewUnauthorizedError("user doesn't login")
		//errJsonReturn(w, r, e)
		return
	}

	vars := mux.Vars(r)
	ns := vars["namespace"]
	if len(ns) == 0 {
		err = errjson.NewErrForbidden("invalid namespace")
		return
	}

	log.Logger.Info(user + " get " + ns + " ugroup info")

	nsJson, err := globalClient.GetNsUgroup(ns)
	if err != nil {
		err = checkDbErr(err)
		return
	}
	fmt.Fprintf(w, string(nsJson))
	return
}

func addUgroup(w http.ResponseWriter, r *http.Request) (err error) {
	user, err := getRequestUser(w, r)
	if err != nil {
		err = errjson.NewUnauthorizedError("user doesn't login")
		//errJsonReturn(w, r, e)
		return
	}

	log.Logger.Info(user + " add  new  ugroup")
	decoder := json.NewDecoder(r.Body)
	var ug database.UserGroup
	err = decoder.Decode(&ug)
	if err != nil {
		panic(err)
	}

	nsJson, err := globalClient.AddUgroup(ug)
	if err != nil {
		err = checkDbErr(err)
		return
	}
	fmt.Fprintf(w, string(nsJson))
	return
}
