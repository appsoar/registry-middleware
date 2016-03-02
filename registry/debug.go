package debug

import (
	"fmt"
)

var (
	debug = false
)

func print(dat interface{}) {
	if debug {
		fmt.Println(dat)
	}

}

func init() {
	debug = true
}
