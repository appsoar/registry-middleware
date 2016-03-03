package debug

import (
	"fmt"
)

var (
	debug = false
)

func Print(dat ...interface{}) {
	if debug {
		for _, v := range dat {
			fmt.Println(v)
		}
	}

}

func init() {
	debug = true
}
