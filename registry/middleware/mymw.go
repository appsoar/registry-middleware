package middleware

import (
	"fmt"
	"net/http"
)

type mymw struct {
	test string
}

func NewMymw() *mymw {
	return &mymw{"test"}
}

func (my *mymw) ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	next(rw, r)
	for k, v := range rw.Header() {
		fmt.Printf("%v:%v\n", k, v)
	}
}
