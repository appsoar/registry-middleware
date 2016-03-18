package handler

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"scheduler/errjson"
	"scheduler/log"
)

func getRepos(w http.ResponseWriter, r *http.Request) (err error) {
	user, err := getRequestUser(w, r)
	if err != nil {
		err = errjson.NewUnauthorizedError("user doesn't login")
		//errJsonReturn(w, r, e)
		return
	}

	log.Logger.Info(user + " get repositories")

	nsJson, err := globalClient.GetRepositories()
	if err != nil {
		err = errjson.NewInternalServerError("can't get ns info")
		return
	}
	fmt.Fprintf(w, string(nsJson))
	return
}

func getNsRepos(w http.ResponseWriter, r *http.Request) (err error) {
	user, err := getRequestUser(w, r)
	if err != nil {
		err = errjson.NewUnauthorizedError("user doesn't login")
		//errJsonReturn(w, r, e)
		return
	}

	log.Logger.Info(user + " get repositories")

	vars := mux.Vars(r)
	ns := vars["namespace"]

	if len(ns) == 0 {
		err = errjson.NewNotValidEntityError("invalid namespace")
		return
	}

	nsJson, err := globalClient.GetNsRepos(ns)
	if err != nil {
		err = errjson.NewInternalServerError("can't get ns info")
		return
	}
	fmt.Fprintf(w, string(nsJson))
	return
}

func getUserRepos(w http.ResponseWriter, r *http.Request) (err error) {
	user, err := getRequestUser(w, r)
	if err != nil {
		err = errjson.NewUnauthorizedError("user doesn't login")
		//errJsonReturn(w, r, e)
		return
	}

	log.Logger.Info(user + " get repositories")

	vars := mux.Vars(r)
	ns := vars["user"]

	if len(ns) == 0 {
		err = errjson.NewNotValidEntityError("invalid namespace")
		return
	}

	nsJson, err := globalClient.GetUserRepos(ns)
	if err != nil {
		err = errjson.NewInternalServerError("can't get ns info")
		return
	}
	fmt.Fprintf(w, string(nsJson))
	return
}

func listRepoTags(w http.ResponseWriter, r *http.Request) (err error) {
	user, err := getRequestUser(w, r)
	if err != nil {
		err = errjson.NewUnauthorizedError("user doesn't login")
		//errJsonReturn(w, r, e)
		return
	}

	log.Logger.Info(user + " get repositories")

	vars := mux.Vars(r)
	name := vars["usernameOrNamespace"]
	repoName := vars["repoName"]

	if len(repoName) == 0 {
		err = errjson.NewNotValidEntityError("invalid namespace")
		return
	}

	nsJson, err := globalClient.ListRepoTags(name, repoName)
	if err != nil {
		err = errjson.NewInternalServerError("can't get repo info")
		return
	}
	fmt.Fprintf(w, string(nsJson))
	return
}

func getTagImage(w http.ResponseWriter, r *http.Request) (err error) {
	user, err := getRequestUser(w, r)
	if err != nil {
		err = errjson.NewUnauthorizedError("user doesn't login")
		//errJsonReturn(w, r, e)
		return
	}

	log.Logger.Info(user + " get repositories")

	vars := mux.Vars(r)
	name := vars["usernameOrNamespace"]
	repoName := vars["repoName"]
	tagName := vars["tagName"]

	if len(repoName) == 0 || len(tagName) == 0 {
		err = errjson.NewNotValidEntityError("invalid namespace")
		return
	}

	nsJson, err := globalClient.GetTagImage(name, repoName, tagName)
	if err != nil {
		err = errjson.NewInternalServerError("can't get repo info")
		return
	}
	fmt.Fprintf(w, string(nsJson))
	return
}
