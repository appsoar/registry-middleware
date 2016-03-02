package main

import (
	//	"fmt"
	"net/http"
	"registry/server"
)

type ServerOpts struct {
	host string
}

func main() {
	opts := ServerOpts{host: "192.168.2.110:9090"}
	router := server.NewRouter()
	http.ListenAndServe(opts.host, router)

}
