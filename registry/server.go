package server

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

type ServerOpts struct {
	host string
}

func notfound(w http.ResponseWriter, r *http.Request) {
	//do something
}

func ServerStart(opt ServerOpts) {
	r := mux.NewRouter()
	r.NotFoundHandler = http.HandlerFunc(notfound)
	//r.HandlerFunc("/", HomeHandler).Methods("GET")
	//r.HandlerFunc("/foo", FooHandler).Methods("POST")
	//r.HandlerFunc("/foo/{filename}",FileHandler).Methods("POST")
	err := http.ListenAndServe(opt.host, r)
	if err != nil {
		//do something
	}
}
