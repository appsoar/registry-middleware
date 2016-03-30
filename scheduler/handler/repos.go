package handler

import (
	"fmt"
	"net/http"
	"scheduler/Godeps/_workspace/src/github.com/gorilla/mux"
	"scheduler/errjson"
	"scheduler/log"
)

func getRepos(w http.ResponseWriter, r *http.Request) (err error) {
	user, err := getRequestUser(w, r)
	if err != nil {
		err = errjson.NewUnauthorizedError("user doesn't login")
		return
	}

	log.Logger.Info(r.RemoteAddr + ":" + user + " get repositories")

	nsJson, err := globalClient.GetRepositories()
	if err != nil {
		err = checkDbErr(err)
		return
	}
	fmt.Fprintf(w, string(nsJson))
	return
}

func getNsRepos(w http.ResponseWriter, r *http.Request) (err error) {
	user, err := getRequestUser(w, r)
	if err != nil {
		err = errjson.NewUnauthorizedError("user doesn't login")
		return
	}

	log.Logger.Info(r.RemoteAddr + ":" + user + " get repositories")

	vars := mux.Vars(r)
	ns := vars["namespace"]

	if len(ns) == 0 {
		err = errjson.NewErrForbidden("invalid namespace")
		return
	}

	nsJson, err := globalClient.GetNsRepos(ns)
	if err != nil {
		err = checkDbErr(err)
		return
	}
	fmt.Fprintf(w, string(nsJson))
	return
}

func getUserRepos(w http.ResponseWriter, r *http.Request) (err error) {
	user, err := getRequestUser(w, r)
	if err != nil {
		err = errjson.NewUnauthorizedError("user doesn't login")
		return
	}

	log.Logger.Info(r.RemoteAddr + ":" + user + " get repositories")

	vars := mux.Vars(r)
	ns := vars["user"]

	if len(ns) == 0 {
		err = errjson.NewErrForbidden("invalid namespace")
		return
	}

	nsJson, err := globalClient.GetUserRepos(ns)
	if err != nil {
		err = checkDbErr(err)
		return
	}
	fmt.Fprintf(w, string(nsJson))
	return
}

func listRepoTags(w http.ResponseWriter, r *http.Request) (err error) {
	user, err := getRequestUser(w, r)
	if err != nil {
		err = errjson.NewUnauthorizedError("user doesn't login")
		return
	}

	log.Logger.Info(r.RemoteAddr + ":" + user + " get repositories")

	vars := mux.Vars(r)
	repoName := vars["repoName"]

	if len(repoName) == 0 {
		err = errjson.NewErrForbidden("invalid repos")
		return
	}

	nsJson, err := globalClient.ListRepoTags(repoName)
	if err != nil {
		err = checkDbErr(err)
		return
	}
	fmt.Fprintf(w, string(nsJson))
	return
}

func getTagImage(w http.ResponseWriter, r *http.Request) (err error) {
	user, err := getRequestUser(w, r)
	if err != nil {
		err = errjson.NewUnauthorizedError("user doesn't login")
		return
	}

	log.Logger.Info(r.RemoteAddr + ":" + user + " get repositories")

	vars := mux.Vars(r)
	repoName := vars["repoName"]
	tagName := vars["tagName"]

	if len(repoName) == 0 || len(tagName) == 0 {
		err = errjson.NewErrForbidden("invalid repo or tag")
		return
	}

	nsJson, err := globalClient.GetTagImage(repoName, tagName)
	if err != nil {
		err = checkDbErr(err)
		return
	}
	fmt.Fprintf(w, string(nsJson))
	return
}
