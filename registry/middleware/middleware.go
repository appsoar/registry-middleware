package middleware

import (
	"github.com/codegangsta/negroni"
)

var negronies = []negroni.Handler{
	negroni.NewRecovery(),
	negroni.NewLogger(),
}

//
func New() *negroni.Negroni {

	n := negroni.New()
	for _, v := range negronies {
		n.Use(v)
	}
	return n
}
