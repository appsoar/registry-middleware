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

	log.Logger.Info(r.RemoteAddr + ":" + user + "get namespace info")
	nsJson, err := globalClient.GetNamespaces()
	if err != nil {
		err = checkDbErr(err)
		return
	}
	fmt.Fprintf(w, string(nsJson))
	return
}

func updateNs(w http.ResponseWriter, r *http.Request) (err error) {
	user, err := getRequestUser(w, r)
	if err != nil {
		err = errjson.NewUnauthorizedError("user doesn't login")
		//errJsonReturn(w, r, e)
		return
	}

	log.Logger.Info(r.RemoteAddr + ":" + user + "update namespace info")
	nsJson, err := globalClient.UpdateNamespace()
	if err != nil {
		err = checkDbErr(err)
		return
	}
	fmt.Fprintf(w, string(nsJson))
	return
}

func deleteNs(w http.ResponseWriter, r *http.Request) (err error) {
	user, err := getRequestUser(w, r)
	if err != nil {
		err = errjson.NewUnauthorizedError("user doesn't login")
		//errJsonReturn(w, r, e)
		return
	}

	vars := mux.Vars(r)
	ns, ok := vars["namespace_name"]
	if !ok {
		panic("missing namespace_name")
	}

	log.Logger.Info(r.RemoteAddr + ":" + user + "delete namespace")
	nsJson, err := globalClient.DeleteNamespace(ns)
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

	log.Logger.Info(r.RemoteAddr + ":" + user + " get " + ns + " namespace info")

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

	log.Logger.Info(r.RemoteAddr + ":" + user + " get " + ns + " ugroup info")

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

	log.Logger.Info(r.RemoteAddr + ":" + user + " add  new  ugroup")
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

func updateUgroup(w http.ResponseWriter, r *http.Request) (err error) {
	user, err := getRequestUser(w, r)
	if err != nil {
		err = errjson.NewUnauthorizedError("user doesn't login")
		//errJsonReturn(w, r, e)
		return
	}

	log.Logger.Info(r.RemoteAddr + ":" + user + "get usergroup")
	nsJson, err := globalClient.UpdateUgroup()
	if err != nil {
		err = checkDbErr(err)
		return
	}
	fmt.Fprintf(w, string(nsJson))
	return
}

func getUgroup(w http.ResponseWriter, r *http.Request) (err error) {
	user, err := getRequestUser(w, r)
	if err != nil {
		err = errjson.NewUnauthorizedError("user doesn't login")
		//errJsonReturn(w, r, e)
		return
	}
	log.Logger.Info(r.RemoteAddr + ":" + user + "get usergroup")

	vars := mux.Vars(r)
	gid, ok := vars["group_id"]
	if !ok {
		panic("group_id missing")
	}

	nsJson, err := globalClient.GetUgroup(gid)
	if err != nil {
		err = checkDbErr(err)
		return
	}
	fmt.Fprintf(w, string(nsJson))
	return
}

func deleteUgroup(w http.ResponseWriter, r *http.Request) (err error) {
	user, err := getRequestUser(w, r)
	if err != nil {
		err = errjson.NewUnauthorizedError("user doesn't login")
		//errJsonReturn(w, r, e)
		return
	}
	log.Logger.Info(r.RemoteAddr + ":" + user + "get usergroup")

	vars := mux.Vars(r)
	gid, ok := vars["group_id"]
	if !ok {
		panic("group_id missing")
	}

	nsJson, err := globalClient.DeleteUgroup(gid)
	if err != nil {
		err = checkDbErr(err)
		return
	}
	fmt.Fprintf(w, string(nsJson))
	return
}
