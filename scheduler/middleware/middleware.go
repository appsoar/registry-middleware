package middleware

import (
	"scheduler/Godeps/_workspace/src/github.com/codegangsta/negroni"
)

var negronies = []negroni.Handler{
	negroni.NewRecovery(),
	//	NewMymw(),
	//	negroni.NewLogger(),
}

//
func New() *negroni.Negroni {

	n := negroni.New()
	for _, v := range negronies {
		n.Use(v)
	}
	return n
}
